package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrSectionNotFound = errors.New("section not found")

type SectionUseCase struct {
	sections repository.SectionRepository
	courses  repository.CourseRepository
}

func NewSectionUseCase(sections repository.SectionRepository, courses repository.CourseRepository) *SectionUseCase {
	return &SectionUseCase{sections: sections, courses: courses}
}

func (uc *SectionUseCase) ListByCourse(courseID string) ([]domain.Section, error) {
	if _, ok, err := uc.courses.FindByID(courseID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrCourseNotFound
	}
	return uc.sections.ListByCourseID(courseID)
}

func (uc *SectionUseCase) Get(id string) (domain.Section, error) {
	s, ok, err := uc.sections.FindByID(id)
	if err != nil {
		return domain.Section{}, err
	}
	if !ok {
		return domain.Section{}, ErrSectionNotFound
	}
	return s, nil
}
