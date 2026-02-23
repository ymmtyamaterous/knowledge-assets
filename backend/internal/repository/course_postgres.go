package repository

import (
	"context"
	"errors"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCourseRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCourseRepository(db *pgxpool.Pool) *PostgresCourseRepository {
	return &PostgresCourseRepository{db: db}
}

func (r *PostgresCourseRepository) List() ([]domain.Course, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, title, description, difficulty, estimated_hours, thumbnail_url, "order", created_at, updated_at
		 FROM courses ORDER BY "order" ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Course
	for rows.Next() {
		c, err := scanCourse(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *PostgresCourseRepository) FindByID(id string) (domain.Course, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, title, description, difficulty, estimated_hours, thumbnail_url, "order", created_at, updated_at
		 FROM courses WHERE id = $1`, id)
	var c domain.Course
	err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Difficulty, &c.EstimatedHour, &c.ThumbnailURL, &c.Order, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Course{}, false, nil
		}
		return domain.Course{}, false, err
	}
	return c, true, nil
}

func scanCourse(rows pgx.Rows) (domain.Course, error) {
	var c domain.Course
	err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Difficulty, &c.EstimatedHour, &c.ThumbnailURL, &c.Order, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}
