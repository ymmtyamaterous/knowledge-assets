package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresQuizRepository struct {
	db *pgxpool.Pool
}

func NewPostgresQuizRepository(db *pgxpool.Pool) *PostgresQuizRepository {
	return &PostgresQuizRepository{db: db}
}

func (r *PostgresQuizRepository) FindByLessonID(lessonID string) (domain.Quiz, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, COALESCE(lesson_id,''), COALESCE(section_id,''), is_mock_exam, COALESCE(time_limit_minutes,0), created_at
		 FROM quizzes WHERE lesson_id = $1`, lessonID)
	return scanQuiz(row)
}

func (r *PostgresQuizRepository) FindBySectionID(sectionID string) (domain.Quiz, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, COALESCE(lesson_id,''), COALESCE(section_id,''), is_mock_exam, COALESCE(time_limit_minutes,0), created_at
		 FROM quizzes WHERE section_id = $1 AND lesson_id IS NULL ORDER BY created_at LIMIT 1`, sectionID)
	return scanQuiz(row)
}

func (r *PostgresQuizRepository) FindByID(id string) (domain.Quiz, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, COALESCE(lesson_id,''), COALESCE(section_id,''), is_mock_exam, COALESCE(time_limit_minutes,0), created_at
		 FROM quizzes WHERE id = $1`, id)
	return scanQuiz(row)
}

func (r *PostgresQuizRepository) ListQuestions(quizID string) ([]domain.QuizQuestion, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, quiz_id, question_text, explanation, "order"
		 FROM quiz_questions WHERE quiz_id = $1 ORDER BY "order" ASC`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.QuizQuestion
	for rows.Next() {
		var q domain.QuizQuestion
		if err := rows.Scan(&q.ID, &q.QuizID, &q.QuestionText, &q.Explanation, &q.Order); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, q := range questions {
		choices, err := r.fetchChoices(q.ID)
		if err != nil {
			return nil, err
		}
		questions[i].Choices = choices
	}
	return questions, nil
}

func (r *PostgresQuizRepository) CreateResult(result domain.UserQuizResult) (domain.UserQuizResult, error) {
	now := time.Now().UTC()
	if result.ID == "" {
		result.ID = fmt.Sprintf("quiz_result_%d", now.UnixNano())
	}
	if result.TakenAt.IsZero() {
		result.TakenAt = now
	}

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO user_quiz_results (id, user_id, quiz_id, score, total, taken_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		result.ID, result.UserID, result.QuizID, result.Score, result.Total, result.TakenAt,
	)
	if err != nil {
		return domain.UserQuizResult{}, err
	}
	return result, nil
}

func (r *PostgresQuizRepository) ListResultsByUserID(userID string) ([]domain.UserQuizResult, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, quiz_id, score, total, taken_at
		 FROM user_quiz_results WHERE user_id = $1 ORDER BY taken_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.UserQuizResult
	for rows.Next() {
		var res domain.UserQuizResult
		if err := rows.Scan(&res.ID, &res.UserID, &res.QuizID, &res.Score, &res.Total, &res.TakenAt); err != nil {
			return nil, err
		}
		list = append(list, res)
	}
	return list, rows.Err()
}

func scanQuiz(row pgx.Row) (domain.Quiz, bool, error) {
	var q domain.Quiz
	err := row.Scan(&q.ID, &q.LessonID, &q.SectionID, &q.IsMockExam, &q.TimeLimitMinutes, &q.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Quiz{}, false, nil
		}
		return domain.Quiz{}, false, err
	}
	return q, true, nil
}

func (r *PostgresQuizRepository) fetchChoices(questionID string) ([]domain.QuizChoice, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, question_id, choice_text, is_correct
		 FROM quiz_choices WHERE question_id = $1`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var choices []domain.QuizChoice
	for rows.Next() {
		var c domain.QuizChoice
		if err := rows.Scan(&c.ID, &c.QuestionID, &c.ChoiceText, &c.IsCorrect); err != nil {
			return nil, err
		}
		choices = append(choices, c)
	}
	return choices, rows.Err()
}
