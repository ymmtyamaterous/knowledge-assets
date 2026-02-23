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
	ID         string        `json:"id"`
	Term       string        `json:"term"`
	Reading    string        `json:"reading"`
	Definition string        `json:"definition"`
	Tags       []GlossaryTag `json:"tags"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

type GlossaryTag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Quiz struct {
	ID               string    `json:"id"`
	LessonID         string    `json:"lessonId"`
	SectionID        string    `json:"sectionId"`
	IsMockExam       bool      `json:"isMockExam"`
	TimeLimitMinutes int       `json:"timeLimitMinutes"`
	CreatedAt        time.Time `json:"createdAt"`
}

type QuizQuestion struct {
	ID           string       `json:"id"`
	QuizID       string       `json:"quizId"`
	QuestionText string       `json:"questionText"`
	Explanation  string       `json:"explanation"`
	Order        int          `json:"order"`
	Choices      []QuizChoice `json:"choices"`
}

type QuizChoice struct {
	ID         string `json:"id"`
	QuestionID string `json:"questionId"`
	ChoiceText string `json:"choiceText"`
	IsCorrect  bool   `json:"isCorrect"`
}

type UserQuizResult struct {
	ID      string    `json:"id"`
	UserID  string    `json:"userId"`
	QuizID  string    `json:"quizId"`
	Score   int       `json:"score"`
	Total   int       `json:"total"`
	TakenAt time.Time `json:"takenAt"`
}

type SectionProgress struct {
	SectionID        string `json:"sectionId"`
	SectionTitle     string `json:"sectionTitle"`
	TotalLessons     int    `json:"totalLessons"`
	CompletedLessons int    `json:"completedLessons"`
	ProgressRate     int    `json:"progressRate"`
}

type CourseProgress struct {
	CourseID         string            `json:"courseId"`
	CourseTitle      string            `json:"courseTitle"`
	TotalLessons     int               `json:"totalLessons"`
	CompletedLessons int               `json:"completedLessons"`
	ProgressRate     int               `json:"progressRate"`
	Sections         []SectionProgress `json:"sections"`
}
