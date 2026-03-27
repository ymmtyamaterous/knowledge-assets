package repository

import (
	"context"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSearchRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSearchRepository(db *pgxpool.Pool) *PostgresSearchRepository {
	return &PostgresSearchRepository{db: db}
}

func (r *PostgresSearchRepository) SearchLessons(query string) ([]domain.SearchLesson, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, title, section_id FROM lessons
		 WHERE title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%'
		 ORDER BY title
		 LIMIT 20`,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.SearchLesson
	for rows.Next() {
		var l domain.SearchLesson
		if err := rows.Scan(&l.ID, &l.Title, &l.SectionID); err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, nil
}

func (r *PostgresSearchRepository) SearchTerms(query string) ([]domain.SearchTerm, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, term, reading FROM glossary_terms
		 WHERE term ILIKE '%' || $1 || '%' OR reading ILIKE '%' || $1 || '%'
		   OR definition ILIKE '%' || $1 || '%'
		 ORDER BY term
		 LIMIT 20`,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.SearchTerm
	for rows.Next() {
		var t domain.SearchTerm
		if err := rows.Scan(&t.ID, &t.Term, &t.Reading); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}
