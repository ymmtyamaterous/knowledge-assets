package usecase

import (
	"testing"

	"asenare/backend/internal/repository"
)

func TestQuizUseCase_Submit(t *testing.T) {
	quizRepo := repository.NewInMemoryQuizRepository()
	lessonRepo := repository.NewInMemoryLessonRepository()
	uc := NewQuizUseCase(quizRepo, lessonRepo)

	detail, err := uc.Get("quiz-fp3-s1-l1")
	if err != nil {
		t.Fatalf("get quiz detail error: %v", err)
	}
	if len(detail.Questions) == 0 {
		t.Fatal("questions must not be empty")
	}

	answers := []QuizAnswer{
		{QuestionID: detail.Questions[0].ID, ChoiceID: detail.Questions[0].Choices[1].ID},
		{QuestionID: detail.Questions[1].ID, ChoiceID: detail.Questions[1].Choices[1].ID},
	}

	result, err := uc.Submit("u1", "quiz-fp3-s1-l1", answers)
	if err != nil {
		t.Fatalf("submit error: %v", err)
	}
	if result.Result.Total != 2 {
		t.Fatalf("unexpected total: %d", result.Result.Total)
	}
	if result.Result.Score != 2 {
		t.Fatalf("unexpected score: %d", result.Result.Score)
	}
}
