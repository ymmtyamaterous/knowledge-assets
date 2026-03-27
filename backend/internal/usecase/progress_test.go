package usecase

import (
	"testing"

	"asenare/backend/internal/repository"
)

func TestProgressUseCase_GetCourseProgress(t *testing.T) {
	progressRepo := repository.NewInMemoryProgressRepository()
	lessonRepo := repository.NewInMemoryLessonRepository()
	courseRepo := repository.NewInMemoryCourseRepository()
	sectionRepo := repository.NewInMemorySectionRepository()
	quizRepo := repository.NewInMemoryQuizRepository()
	noteRepo := repository.NewMemoryNoteRepository()

	uc := NewProgressUseCase(progressRepo, lessonRepo, courseRepo, sectionRepo, quizRepo, noteRepo)

	if _, err := uc.CompleteLesson("u1", "fp3-s1-l1"); err != nil {
		t.Fatalf("complete lesson 1: %v", err)
	}
	if _, err := uc.CompleteLesson("u1", "fp3-s1-l2"); err != nil {
		t.Fatalf("complete lesson 2: %v", err)
	}

	list, err := uc.GetCourseProgress("u1")
	if err != nil {
		t.Fatalf("get course progress error: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("course progress should not be empty")
	}

	var fp3Found bool
	for _, cp := range list {
		if cp.CourseID == "fp3" {
			fp3Found = true
			if cp.CompletedLessons < 2 {
				t.Fatalf("expected completed lessons >= 2, got %d", cp.CompletedLessons)
			}
		}
	}

	if !fp3Found {
		t.Fatal("fp3 progress not found")
	}
}

func TestProgressUseCase_UncompleteLesson(t *testing.T) {
	progressRepo := repository.NewInMemoryProgressRepository()
	lessonRepo := repository.NewInMemoryLessonRepository()
	courseRepo := repository.NewInMemoryCourseRepository()
	sectionRepo := repository.NewInMemorySectionRepository()
	quizRepo := repository.NewInMemoryQuizRepository()
	noteRepo := repository.NewMemoryNoteRepository()

	uc := NewProgressUseCase(progressRepo, lessonRepo, courseRepo, sectionRepo, quizRepo, noteRepo)

	if _, err := uc.CompleteLesson("u1", "fp3-s1-l1"); err != nil {
		t.Fatalf("complete lesson error: %v", err)
	}

	if err := uc.UncompleteLesson("u1", "fp3-s1-l1"); err != nil {
		t.Fatalf("uncomplete lesson error: %v", err)
	}

	progress, err := uc.GetUserProgress("u1")
	if err != nil {
		t.Fatalf("get user progress error: %v", err)
	}
	if len(progress) != 0 {
		t.Fatalf("expected progress to be empty, got %d", len(progress))
	}
}
