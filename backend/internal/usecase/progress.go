package usecase

import (
	"errors"
	"sort"
	"time"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrAlreadyCompleted = errors.New("lesson already completed")

type ProgressUseCase struct {
	progress repository.ProgressRepository
	lessons  repository.LessonRepository
	courses  repository.CourseRepository
	sections repository.SectionRepository
	quizzes  repository.QuizRepository
	notes    repository.NoteRepository
}

func NewProgressUseCase(
	progress repository.ProgressRepository,
	lessons repository.LessonRepository,
	courses repository.CourseRepository,
	sections repository.SectionRepository,
	quizzes repository.QuizRepository,
	notes repository.NoteRepository,
) *ProgressUseCase {
	return &ProgressUseCase{progress: progress, lessons: lessons, courses: courses, sections: sections, quizzes: quizzes, notes: notes}
}

func (uc *ProgressUseCase) CompleteLesson(userID, lessonID string) (domain.UserLessonProgress, error) {
	if _, ok, err := uc.lessons.FindByID(lessonID); err != nil {
		return domain.UserLessonProgress{}, err
	} else if !ok {
		return domain.UserLessonProgress{}, ErrLessonNotFound
	}

	if existing, ok, err := uc.progress.FindByUserAndLesson(userID, lessonID); err != nil {
		return domain.UserLessonProgress{}, err
	} else if ok {
		return existing, nil // 冪等: すでに完了済みならそのまま返す
	}

	return uc.progress.Create(domain.UserLessonProgress{
		UserID:   userID,
		LessonID: lessonID,
	})
}

func (uc *ProgressUseCase) UncompleteLesson(userID, lessonID string) error {
	if _, ok, err := uc.lessons.FindByID(lessonID); err != nil {
		return err
	} else if !ok {
		return ErrLessonNotFound
	}

	return uc.progress.DeleteByUserAndLesson(userID, lessonID)
}

func (uc *ProgressUseCase) GetUserProgress(userID string) ([]domain.UserLessonProgress, error) {
	return uc.progress.ListByUserID(userID)
}

func (uc *ProgressUseCase) GetCourseProgress(userID string) ([]domain.CourseProgress, error) {
	courses, err := uc.courses.List()
	if err != nil {
		return nil, err
	}

	result := make([]domain.CourseProgress, 0, len(courses))
	for _, course := range courses {
		sections, err := uc.sections.ListByCourseID(course.ID)
		if err != nil {
			return nil, err
		}

		courseSummary := domain.CourseProgress{
			CourseID:    course.ID,
			CourseTitle: course.Title,
			Sections:    make([]domain.SectionProgress, 0, len(sections)),
		}

		for _, section := range sections {
			lessons, err := uc.lessons.ListBySectionID(section.ID)
			if err != nil {
				return nil, err
			}

			lessonIDs := make([]string, 0, len(lessons))
			for _, lesson := range lessons {
				lessonIDs = append(lessonIDs, lesson.ID)
			}

			completedCount, err := uc.progress.CountByUserAndCourse(userID, course.ID, lessonIDs)
			if err != nil {
				return nil, err
			}

			sectionTotal := len(lessonIDs)
			sectionRate := 0
			if sectionTotal > 0 {
				sectionRate = (completedCount * 100) / sectionTotal
			}

			courseSummary.Sections = append(courseSummary.Sections, domain.SectionProgress{
				SectionID:        section.ID,
				SectionTitle:     section.Title,
				TotalLessons:     sectionTotal,
				CompletedLessons: completedCount,
				ProgressRate:     sectionRate,
			})

			courseSummary.TotalLessons += sectionTotal
			courseSummary.CompletedLessons += completedCount
		}

		if courseSummary.TotalLessons > 0 {
			courseSummary.ProgressRate = (courseSummary.CompletedLessons * 100) / courseSummary.TotalLessons
		}

		result = append(result, courseSummary)
	}

	return result, nil
}

// GetStreak は今日から遡って連続した学習日数と最長連続日数を返します。
func (uc *ProgressUseCase) GetStreak(userID string) (domain.UserStreak, error) {
	list, err := uc.progress.ListByUserID(userID)
	if err != nil {
		return domain.UserStreak{}, err
	}
	if len(list) == 0 {
		return domain.UserStreak{}, nil
	}

	// 日付（YYYY-MM-DD）をユニーク化してソート
	daySet := make(map[string]struct{}, len(list))
	for _, p := range list {
		key := p.CompletedAt.UTC().Format("2006-01-02")
		daySet[key] = struct{}{}
	}
	days := make([]string, 0, len(daySet))
	for d := range daySet {
		days = append(days, d)
	}
	sort.Strings(days) // 昇順

	// 最終学習日
	lastStudiedAt := days[len(days)-1]

	// 今日から遡って連続日数を計算
	today := time.Now().UTC().Format("2006-01-02")
	current := 0
	check := today
	for {
		if _, ok := daySet[check]; ok {
			current++
			t, _ := time.Parse("2006-01-02", check)
			check = t.AddDate(0, 0, -1).Format("2006-01-02")
		} else {
			// 今日学習していない場合、昨日から遡る
			if check == today && current == 0 {
				yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")
				if _, ok := daySet[yesterday]; ok {
					current++
					t, _ := time.Parse("2006-01-02", yesterday)
					check = t.AddDate(0, 0, -1).Format("2006-01-02")
					continue
				}
			}
			break
		}
	}

	// 全期間で最長連続日数を計算
	longest := 0
	streak := 0
	for i, d := range days {
		if i == 0 {
			streak = 1
		} else {
			prev, _ := time.Parse("2006-01-02", days[i-1])
			cur, _ := time.Parse("2006-01-02", d)
			diff := cur.Sub(prev).Hours() / 24
			if diff == 1 {
				streak++
			} else {
				streak = 1
			}
		}
		if streak > longest {
			longest = streak
		}
	}

	return domain.UserStreak{
		CurrentStreak: current,
		LongestStreak: longest,
		LastStudiedAt: lastStudiedAt,
	}, nil
}

// GetStats は学習統計サマリを返します。
func (uc *ProgressUseCase) GetStats(userID string) (domain.UserStats, error) {
	// 完了レッスン数・学習日数
	progList, err := uc.progress.ListByUserID(userID)
	if err != nil {
		return domain.UserStats{}, err
	}
	daySet := make(map[string]struct{}, len(progList))
	for _, p := range progList {
		daySet[p.CompletedAt.UTC().Format("2006-01-02")] = struct{}{}
	}

	// クイズ平均正答率
	results, err := uc.quizzes.ListResultsByUserID(userID)
	if err != nil {
		return domain.UserStats{}, err
	}
	var avgScore float64
	if len(results) > 0 {
		var sum float64
		for _, r := range results {
			if r.Total > 0 {
				sum += float64(r.Score) / float64(r.Total) * 100
			}
		}
		avgScore = sum / float64(len(results))
	}

	// メモ件数
	notes, err := uc.notes.ListByUserID(userID)
	if err != nil {
		return domain.UserStats{}, err
	}

	return domain.UserStats{
		TotalCompletedLessons: len(progList),
		TotalStudyDays:        len(daySet),
		TotalNotes:            len(notes),
		AverageQuizScore:      avgScore,
	}, nil
}

// GetCalendar は指定年の日別完了レッスン数を返します（year=0 の場合は今年）。
func (uc *ProgressUseCase) GetCalendar(userID string, year int) (domain.UserCalendar, error) {
	if year == 0 {
		year = time.Now().UTC().Year()
	}
	list, err := uc.progress.ListByUserID(userID)
	if err != nil {
		return domain.UserCalendar{}, err
	}
	counts := make(map[string]int)
	for _, p := range list {
		if p.CompletedAt.UTC().Year() != year {
			continue
		}
		key := p.CompletedAt.UTC().Format("2006-01-02")
		counts[key]++
	}
	days := make([]domain.CalendarDay, 0, len(counts))
	for date, count := range counts {
		days = append(days, domain.CalendarDay{Date: date, Count: count})
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Date < days[j].Date })
	return domain.UserCalendar{Days: days}, nil
}
