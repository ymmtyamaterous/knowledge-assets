package repository

import (
	"context"
	"fmt"
	"time"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresBadgeRepository struct {
	db *pgxpool.Pool
}

func NewPostgresBadgeRepository(db *pgxpool.Pool) *PostgresBadgeRepository {
	return &PostgresBadgeRepository{db: db}
}

func (r *PostgresBadgeRepository) FindByCondition(conditionType, conditionID string) (domain.Badge, bool, error) {
	ctx := context.Background()
	row := r.db.QueryRow(ctx,
		`SELECT id, name, description, image_url, condition_type, condition_id
		 FROM badges WHERE condition_type = $1 AND condition_id = $2`,
		conditionType, conditionID,
	)

	var b domain.Badge
	if err := row.Scan(&b.ID, &b.Name, &b.Description, &b.ImageURL, &b.ConditionType, &b.ConditionID); err != nil {
		return domain.Badge{}, false, nil
	}
	return b, true, nil
}

func (r *PostgresBadgeRepository) CreateUserBadge(userID, badgeID string) (domain.UserBadge, error) {
	ctx := context.Background()
	id := fmt.Sprintf("ub-%s-%s-%d", userID, badgeID, time.Now().UnixNano())
	now := time.Now().UTC()

	_, err := r.db.Exec(ctx,
		`INSERT INTO user_badges (id, user_id, badge_id, earned_at) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (user_id, badge_id) DO NOTHING`,
		id, userID, badgeID, now,
	)
	if err != nil {
		return domain.UserBadge{}, err
	}

	// 既存レコードも含めて取得（ON CONFLICT DO NOTHING の場合に対応）
	return r.fetchUserBadge(ctx, userID, badgeID)
}

func (r *PostgresBadgeRepository) fetchUserBadge(ctx context.Context, userID, badgeID string) (domain.UserBadge, error) {
	row := r.db.QueryRow(ctx,
		`SELECT ub.id, ub.user_id, ub.earned_at,
		        b.id, b.name, b.description, b.image_url, b.condition_type, b.condition_id
		 FROM user_badges ub
		 JOIN badges b ON b.id = ub.badge_id
		 WHERE ub.user_id = $1 AND ub.badge_id = $2`,
		userID, badgeID,
	)

	var ub domain.UserBadge
	if err := row.Scan(
		&ub.ID, &ub.UserID, &ub.EarnedAt,
		&ub.Badge.ID, &ub.Badge.Name, &ub.Badge.Description, &ub.Badge.ImageURL,
		&ub.Badge.ConditionType, &ub.Badge.ConditionID,
	); err != nil {
		return domain.UserBadge{}, err
	}
	return ub, nil
}

func (r *PostgresBadgeRepository) ListByUserID(userID string) ([]domain.UserBadge, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT ub.id, ub.user_id, ub.earned_at,
		        b.id, b.name, b.description, b.image_url, b.condition_type, b.condition_id
		 FROM user_badges ub
		 JOIN badges b ON b.id = ub.badge_id
		 WHERE ub.user_id = $1
		 ORDER BY ub.earned_at ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.UserBadge
	for rows.Next() {
		var ub domain.UserBadge
		if err := rows.Scan(
			&ub.ID, &ub.UserID, &ub.EarnedAt,
			&ub.Badge.ID, &ub.Badge.Name, &ub.Badge.Description, &ub.Badge.ImageURL,
			&ub.Badge.ConditionType, &ub.Badge.ConditionID,
		); err != nil {
			return nil, err
		}
		result = append(result, ub)
	}
	return result, nil
}

func (r *PostgresBadgeRepository) ExistsUserBadge(userID, badgeID string) (bool, error) {
	ctx := context.Background()
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM user_badges WHERE user_id = $1 AND badge_id = $2)`,
		userID, badgeID,
	).Scan(&exists)
	return exists, err
}
