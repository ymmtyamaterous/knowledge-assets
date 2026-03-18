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

type PostgresNoteRepository struct {
	db *pgxpool.Pool
}

func NewPostgresNoteRepository(db *pgxpool.Pool) *PostgresNoteRepository {
	return &PostgresNoteRepository{db: db}
}

func (r *PostgresNoteRepository) FindByUserAndLesson(userID, lessonID string) (domain.UserNote, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, lesson_id, content, created_at, updated_at
		 FROM user_notes WHERE user_id = $1 AND lesson_id = $2`, userID, lessonID)
	var n domain.UserNote
	err := row.Scan(&n.ID, &n.UserID, &n.LessonID, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserNote{}, false, nil
		}
		return domain.UserNote{}, false, err
	}
	return n, true, nil
}

func (r *PostgresNoteRepository) Upsert(note domain.UserNote) (domain.UserNote, error) {
	now := time.Now().UTC()
	if note.ID == "" {
		note.ID = fmt.Sprintf("note_%d", now.UnixNano())
	}
	note.UpdatedAt = now

	row := r.db.QueryRow(context.Background(),
		`INSERT INTO user_notes (id, user_id, lesson_id, content, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (user_id, lesson_id) DO UPDATE
		   SET content = EXCLUDED.content, updated_at = EXCLUDED.updated_at
		 RETURNING id, user_id, lesson_id, content, created_at, updated_at`,
		note.ID, note.UserID, note.LessonID, note.Content, now, now,
	)
	var n domain.UserNote
	if err := row.Scan(&n.ID, &n.UserID, &n.LessonID, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
		return domain.UserNote{}, err
	}
	return n, nil
}

func (r *PostgresNoteRepository) ListByUserID(userID string) ([]domain.UserNote, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, lesson_id, content, created_at, updated_at
		 FROM user_notes WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.UserNote
	for rows.Next() {
		var n domain.UserNote
		if err := rows.Scan(&n.ID, &n.UserID, &n.LessonID, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}
