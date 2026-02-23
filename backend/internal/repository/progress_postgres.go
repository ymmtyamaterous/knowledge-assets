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

type PostgresProgressRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProgressRepository(db *pgxpool.Pool) *PostgresProgressRepository {
	return &PostgresProgressRepository{db: db}
}

func (r *PostgresProgressRepository) FindByUserAndLesson(userID, lessonID string) (domain.UserLessonProgress, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, lesson_id, completed_at
		 FROM user_lesson_progress WHERE user_id = $1 AND lesson_id = $2`, userID, lessonID)
	var p domain.UserLessonProgress
	err := row.Scan(&p.ID, &p.UserID, &p.LessonID, &p.CompletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserLessonProgress{}, false, nil
		}
		return domain.UserLessonProgress{}, false, err
	}
	return p, true, nil
}

func (r *PostgresProgressRepository) Create(p domain.UserLessonProgress) (domain.UserLessonProgress, error) {
	now := time.Now().UTC()
	if p.ID == "" {
		p.ID = fmt.Sprintf("prog_%d", now.UnixNano())
	}
	if p.CompletedAt.IsZero() {
		p.CompletedAt = now
	}

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO user_lesson_progress (id, user_id, lesson_id, completed_at)
		 VALUES ($1, $2, $3, $4)`,
		p.ID, p.UserID, p.LessonID, p.CompletedAt,
	)
	if err != nil {
		return domain.UserLessonProgress{}, err
	}
	return p, nil
}

func (r *PostgresProgressRepository) DeleteByUserAndLesson(userID, lessonID string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM user_lesson_progress WHERE user_id = $1 AND lesson_id = $2`, userID, lessonID)
	return err
}

func (r *PostgresProgressRepository) ListByUserID(userID string) ([]domain.UserLessonProgress, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, lesson_id, completed_at
		 FROM user_lesson_progress WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.UserLessonProgress
	for rows.Next() {
		var p domain.UserLessonProgress
		if err := rows.Scan(&p.ID, &p.UserID, &p.LessonID, &p.CompletedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *PostgresProgressRepository) CountByUserAndCourse(userID, _ string, lessonIDs []string) (int, error) {
	if len(lessonIDs) == 0 {
		return 0, nil
	}

	rows, err := r.db.Query(context.Background(),
		`SELECT COUNT(*) FROM user_lesson_progress
		 WHERE user_id = $1 AND lesson_id = ANY($2)`,
		userID, lessonIDs,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, rows.Err()
}
