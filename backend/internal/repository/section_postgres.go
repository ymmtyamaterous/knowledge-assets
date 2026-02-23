package repository

import (
	"context"
	"errors"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSectionRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSectionRepository(db *pgxpool.Pool) *PostgresSectionRepository {
	return &PostgresSectionRepository{db: db}
}

func (r *PostgresSectionRepository) ListByCourseID(courseID string) ([]domain.Section, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, course_id, title, description, "order", created_at, updated_at
		 FROM sections WHERE course_id = $1 ORDER BY "order" ASC`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Section
	for rows.Next() {
		var s domain.Section
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.Order, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *PostgresSectionRepository) FindByID(id string) (domain.Section, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, course_id, title, description, "order", created_at, updated_at
		 FROM sections WHERE id = $1`, id)
	var s domain.Section
	err := row.Scan(&s.ID, &s.CourseID, &s.Title, &s.Description, &s.Order, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Section{}, false, nil
		}
		return domain.Section{}, false, err
	}
	return s, true, nil
}
