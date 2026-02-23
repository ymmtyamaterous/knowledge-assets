package repository

import (
	"context"
	"errors"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresLessonRepository struct {
	db *pgxpool.Pool
}

func NewPostgresLessonRepository(db *pgxpool.Pool) *PostgresLessonRepository {
	return &PostgresLessonRepository{db: db}
}

func (r *PostgresLessonRepository) ListBySectionID(sectionID string) ([]domain.Lesson, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, section_id, title, content, "order", created_at, updated_at
		 FROM lessons WHERE section_id = $1 ORDER BY "order" ASC`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Lesson
	for rows.Next() {
		var l domain.Lesson
		if err := rows.Scan(&l.ID, &l.SectionID, &l.Title, &l.Content, &l.Order, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

func (r *PostgresLessonRepository) FindByID(id string) (domain.Lesson, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, section_id, title, content, "order", created_at, updated_at
		 FROM lessons WHERE id = $1`, id)
	var l domain.Lesson
	err := row.Scan(&l.ID, &l.SectionID, &l.Title, &l.Content, &l.Order, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Lesson{}, false, nil
		}
		return domain.Lesson{}, false, err
	}
	return l, true, nil
}
