package repository

import "asenare/backend/internal/domain"

type UserRepository interface {
	Create(user domain.User) (domain.User, error)
	FindByEmail(email string) (domain.User, bool, error)
	FindByID(id string) (domain.User, bool, error)
	Update(user domain.User) (domain.User, error)
}

type CourseRepository interface {
	List() ([]domain.Course, error)
	FindByID(id string) (domain.Course, bool, error)
}

type SectionRepository interface {
	ListByCourseID(courseID string) ([]domain.Section, error)
	FindByID(id string) (domain.Section, bool, error)
}

type LessonRepository interface {
	ListBySectionID(sectionID string) ([]domain.Lesson, error)
	FindByID(id string) (domain.Lesson, bool, error)
}

type ProgressRepository interface {
	FindByUserAndLesson(userID, lessonID string) (domain.UserLessonProgress, bool, error)
	Create(p domain.UserLessonProgress) (domain.UserLessonProgress, error)
	DeleteByUserAndLesson(userID, lessonID string) error
	ListByUserID(userID string) ([]domain.UserLessonProgress, error)
	CountByUserAndCourse(userID, courseID string, lessonIDs []string) (int, error)
}

type GlossaryRepository interface {
	List(tagID string) ([]domain.GlossaryTerm, error)
	FindByID(id string) (domain.GlossaryTerm, bool, error)
	ListTags() ([]domain.GlossaryTag, error)
}

type QuizRepository interface {
	FindByLessonID(lessonID string) (domain.Quiz, bool, error)
	FindBySectionID(sectionID string) (domain.Quiz, bool, error)
	FindByID(id string) (domain.Quiz, bool, error)
	ListQuestions(quizID string) ([]domain.QuizQuestion, error)
	CreateResult(result domain.UserQuizResult) (domain.UserQuizResult, error)
	ListResultsByUserID(userID string) ([]domain.UserQuizResult, error)
}

type NoteRepository interface {
	FindByUserAndLesson(userID, lessonID string) (domain.UserNote, bool, error)
	Upsert(note domain.UserNote) (domain.UserNote, error)
	ListByUserID(userID string) ([]domain.UserNote, error)
}

type BadgeRepository interface {
	FindByCondition(conditionType, conditionID string) (domain.Badge, bool, error)
	CreateUserBadge(userID, badgeID string) (domain.UserBadge, error)
	ListByUserID(userID string) ([]domain.UserBadge, error)
	ExistsUserBadge(userID, badgeID string) (bool, error)
}

type SearchRepository interface {
	SearchLessons(query string) ([]domain.SearchLesson, error)
	SearchTerms(query string) ([]domain.SearchTerm, error)
}
