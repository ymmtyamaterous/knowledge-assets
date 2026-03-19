package repository

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryQuizRepository struct {
	mu        sync.RWMutex
	quizzes   map[string]domain.Quiz
	questions map[string][]domain.QuizQuestion
	results   []domain.UserQuizResult
}

func NewInMemoryQuizRepository() *InMemoryQuizRepository {
	now := time.Now().UTC()
	lessonQuiz := domain.Quiz{
		ID:               "quiz-fp3-s1-l1",
		LessonID:         "fp3-s1-l1",
		SectionID:        "fp3-s1",
		IsMockExam:       false,
		TimeLimitMinutes: 10,
		CreatedAt:        now,
	}

	sectionQuiz := domain.Quiz{
		ID:               "quiz-fp3-s1",
		LessonID:         "",
		SectionID:        "fp3-s1",
		IsMockExam:       false,
		TimeLimitMinutes: 15,
		CreatedAt:        now,
	}

	lessonQuestions := []domain.QuizQuestion{
		{
			ID:           "q-fp3-1",
			QuizID:       lessonQuiz.ID,
			QuestionText: "ライフプランの説明として最も適切なものはどれですか？",
			Explanation:  "ライフプランは人生全体の目標を可視化し、必要資金を見積もる長期計画です。",
			Order:        1,
			Choices: []domain.QuizChoice{
				{ID: "q-fp3-1-c1", QuestionID: "q-fp3-1", ChoiceText: "短期の節約術だけをまとめたもの", IsCorrect: false},
				{ID: "q-fp3-1-c2", QuestionID: "q-fp3-1", ChoiceText: "人生のイベントと資金計画を含む長期設計", IsCorrect: true},
				{ID: "q-fp3-1-c3", QuestionID: "q-fp3-1", ChoiceText: "投資商品だけを比較する表", IsCorrect: false},
			},
		},
		{
			ID:           "q-fp3-2",
			QuizID:       lessonQuiz.ID,
			QuestionText: "ライフプラン作成の最初のステップはどれですか？",
			Explanation:  "最初は現状把握です。収入・支出・資産・負債の整理から始めます。",
			Order:        2,
			Choices: []domain.QuizChoice{
				{ID: "q-fp3-2-c1", QuestionID: "q-fp3-2", ChoiceText: "必要資金の算出", IsCorrect: false},
				{ID: "q-fp3-2-c2", QuestionID: "q-fp3-2", ChoiceText: "現状把握", IsCorrect: true},
				{ID: "q-fp3-2-c3", QuestionID: "q-fp3-2", ChoiceText: "資金調達方法の検討", IsCorrect: false},
			},
		},
	}

	sectionQuestions := []domain.QuizQuestion{
		{
			ID:           "q-fp3-s1-1",
			QuizID:       sectionQuiz.ID,
			QuestionText: "キャッシュフロー表に記載する主な項目はどれですか？",
			Explanation:  "キャッシュフロー表には年ごとの収入・支出・年間収支・貯蓄残高を記載します。",
			Order:        1,
			Choices: []domain.QuizChoice{
				{ID: "q-fp3-s1-1-c1", QuestionID: "q-fp3-s1-1", ChoiceText: "収入・支出・年間収支・貯蓄残高", IsCorrect: true},
				{ID: "q-fp3-s1-1-c2", QuestionID: "q-fp3-s1-1", ChoiceText: "株価・金利・為替・物価", IsCorrect: false},
				{ID: "q-fp3-s1-1-c3", QuestionID: "q-fp3-s1-1", ChoiceText: "資産・負債・純資産のみ", IsCorrect: false},
			},
		},
	}

	return &InMemoryQuizRepository{
		quizzes: map[string]domain.Quiz{
			lessonQuiz.ID:  lessonQuiz,
			sectionQuiz.ID: sectionQuiz,
		},
		questions: map[string][]domain.QuizQuestion{
			lessonQuiz.ID:  lessonQuestions,
			sectionQuiz.ID: sectionQuestions,
		},
		results: []domain.UserQuizResult{},
	}
}

func (r *InMemoryQuizRepository) FindByLessonID(lessonID string) (domain.Quiz, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, q := range r.quizzes {
		if q.LessonID == lessonID {
			return q, true, nil
		}
	}
	return domain.Quiz{}, false, nil
}

func (r *InMemoryQuizRepository) FindBySectionID(sectionID string) (domain.Quiz, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, q := range r.quizzes {
		if q.SectionID == sectionID && q.LessonID == "" {
			return q, true, nil
		}
	}
	return domain.Quiz{}, false, nil
}

func (r *InMemoryQuizRepository) FindByID(id string) (domain.Quiz, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	q, ok := r.quizzes[id]
	return q, ok, nil
}

func (r *InMemoryQuizRepository) ListQuestions(quizID string) ([]domain.QuizQuestion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := r.questions[quizID]
	copied := make([]domain.QuizQuestion, len(list))
	copy(copied, list)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Order < copied[j].Order
	})
	return copied, nil
}

func (r *InMemoryQuizRepository) CreateResult(result domain.UserQuizResult) (domain.UserQuizResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	if result.ID == "" {
		result.ID = fmt.Sprintf("quiz_result_%d", now.UnixNano())
	}
	if result.TakenAt.IsZero() {
		result.TakenAt = now
	}
	r.results = append(r.results, result)
	return result, nil
}

func (r *InMemoryQuizRepository) ListResultsByUserID(userID string) ([]domain.UserQuizResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.UserQuizResult, 0)
	for _, result := range r.results {
		if result.UserID == userID {
			list = append(list, result)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].TakenAt.After(list[j].TakenAt)
	})
	return list, nil
}
