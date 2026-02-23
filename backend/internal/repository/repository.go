package repository

import "asenare/backend/internal/domain"

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	FindByEmail(email string) (domain.User, bool, error)
	FindByID(id string) (domain.User, bool, error)
}

type CourseRepository interface {
	List() ([]domain.Course, error)
	FindByID(id string) (domain.Course, bool, error)
}
