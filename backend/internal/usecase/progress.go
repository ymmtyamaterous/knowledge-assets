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
}

func NewProgressUseCase(progress repository.ProgressRepository, lessons repository.LessonRepository) *ProgressUseCase {
	return &ProgressUseCase{progress: progress, lessons: lessons}
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

func (uc *ProgressUseCase) GetUserProgress(userID string) ([]domain.UserLessonProgress, error) {
	return uc.progress.ListByUserID(userID)
}
