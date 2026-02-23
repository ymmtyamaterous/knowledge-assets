package domain

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Username     string    `json:"username"`
	AvatarURL    string    `json:"avatarUrl"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Course struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Difficulty    string    `json:"difficulty"`
	EstimatedHour int       `json:"estimatedHour"`
	ThumbnailURL  string    `json:"thumbnailUrl"`
	Order         int       `json:"order"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Section struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"courseId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Lesson struct {
	ID        string    `json:"id"`
	SectionID string    `json:"sectionId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserLessonProgress struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	LessonID    string    `json:"lessonId"`
	CompletedAt time.Time `json:"completedAt"`
}

type GlossaryTerm struct {
	ID         string    `json:"id"`
	Term       string    `json:"term"`
	Reading    string    `json:"reading"`
	Definition string    `json:"definition"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
