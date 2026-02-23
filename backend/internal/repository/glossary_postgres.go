package repository

import (
	"context"
	"errors"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresGlossaryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresGlossaryRepository(db *pgxpool.Pool) *PostgresGlossaryRepository {
	return &PostgresGlossaryRepository{db: db}
}

func (r *PostgresGlossaryRepository) List(tagID string) ([]domain.GlossaryTerm, error) {
	var rows pgx.Rows
	var err error

	if tagID != "" {
		rows, err = r.db.Query(context.Background(),
			`SELECT DISTINCT gt.id, gt.term, gt.reading, gt.definition, gt.created_at, gt.updated_at
			 FROM glossary_terms gt
			 JOIN glossary_term_tags gtt ON gtt.term_id = gt.id
			 WHERE gtt.tag_id = $1
			 ORDER BY gt.reading ASC`, tagID)
	} else {
		rows, err = r.db.Query(context.Background(),
			`SELECT id, term, reading, definition, created_at, updated_at
			 FROM glossary_terms ORDER BY reading ASC`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []domain.GlossaryTerm
	for rows.Next() {
		var t domain.GlossaryTerm
		if err := rows.Scan(&t.ID, &t.Term, &t.Reading, &t.Definition, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, t := range terms {
		tags, err := r.fetchTagsForTerm(t.ID)
		if err != nil {
			return nil, err
		}
		terms[i].Tags = tags
	}
	return terms, nil
}

func (r *PostgresGlossaryRepository) FindByID(id string) (domain.GlossaryTerm, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, term, reading, definition, created_at, updated_at
		 FROM glossary_terms WHERE id = $1`, id)
	var t domain.GlossaryTerm
	err := row.Scan(&t.ID, &t.Term, &t.Reading, &t.Definition, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.GlossaryTerm{}, false, nil
		}
		return domain.GlossaryTerm{}, false, err
	}

	tags, err := r.fetchTagsForTerm(t.ID)
	if err != nil {
		return domain.GlossaryTerm{}, false, err
	}
	t.Tags = tags
	return t, true, nil
}

func (r *PostgresGlossaryRepository) ListTags() ([]domain.GlossaryTag, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, name, created_at FROM glossary_tags ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.GlossaryTag
	for rows.Next() {
		var tag domain.GlossaryTag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *PostgresGlossaryRepository) fetchTagsForTerm(termID string) ([]domain.GlossaryTag, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT gt.id, gt.name, gt.created_at
		 FROM glossary_tags gt
		 JOIN glossary_term_tags gtt ON gtt.tag_id = gt.id
		 WHERE gtt.term_id = $1
		 ORDER BY gt.name ASC`, termID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.GlossaryTag
	for rows.Next() {
		var tag domain.GlossaryTag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
