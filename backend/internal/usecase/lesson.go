package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrLessonNotFound = errors.New("lesson not found")

type LessonUseCase struct {
	lessons  repository.LessonRepository
	sections repository.SectionRepository
}

func NewLessonUseCase(lessons repository.LessonRepository, sections repository.SectionRepository) *LessonUseCase {
	return &LessonUseCase{lessons: lessons, sections: sections}
}

func (uc *LessonUseCase) ListBySection(sectionID string) ([]domain.Lesson, error) {
	if _, ok, err := uc.sections.FindByID(sectionID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrSectionNotFound
	}
	return uc.lessons.ListBySectionID(sectionID)
}

func (uc *LessonUseCase) Get(id string) (domain.Lesson, error) {
	l, ok, err := uc.lessons.FindByID(id)
	if err != nil {
		return domain.Lesson{}, err
	}
	if !ok {
		return domain.Lesson{}, ErrLessonNotFound
	}
	return l, nil
}
