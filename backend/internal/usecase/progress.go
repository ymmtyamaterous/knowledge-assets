package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrAlreadyCompleted = errors.New("lesson already completed")

type ProgressUseCase struct {
	progress repository.ProgressRepository
	lessons  repository.LessonRepository
	courses  repository.CourseRepository
	sections repository.SectionRepository
}

func NewProgressUseCase(
	progress repository.ProgressRepository,
	lessons repository.LessonRepository,
	courses repository.CourseRepository,
	sections repository.SectionRepository,
) *ProgressUseCase {
	return &ProgressUseCase{progress: progress, lessons: lessons, courses: courses, sections: sections}
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
