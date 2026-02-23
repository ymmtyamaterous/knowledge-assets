package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrCourseNotFound = errors.New("course not found")

type CourseUseCase struct {
	courses repository.CourseRepository
}

func NewCourseUseCase(courses repository.CourseRepository) *CourseUseCase {
	return &CourseUseCase{courses: courses}
}

func (uc *CourseUseCase) List() ([]domain.Course, error) {
	return uc.courses.List()
}

func (uc *CourseUseCase) Get(id string) (domain.Course, error) {
	c, ok, err := uc.courses.FindByID(id)
	if err != nil {
		return domain.Course{}, err
	}
	if !ok {
		return domain.Course{}, ErrCourseNotFound
	}
	return c, nil
}
