package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"asenare/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user domain.User) (domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(user.Email))
	if email == "" {
		return domain.User{}, errors.New("email is required")
	}

	now := time.Now().UTC()
	if user.ID == "" {
		user.ID = "u_" + strings.ReplaceAll(now.Format("20060102150405.000000000"), ".", "")
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now
	user.Email = email

	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, username, avatar_url, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Email, user.PasswordHash, user.Username, user.AvatarURL, string(user.Role), user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return domain.User{}, errors.New("email already exists")
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (domain.User, bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	row := r.db.QueryRow(context.Background(),
		`SELECT id, email, password_hash, username, avatar_url, role, created_at, updated_at
		 FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *PostgresUserRepository) FindByID(id string) (domain.User, bool, error) {
	row := r.db.QueryRow(context.Background(),
		`SELECT id, email, password_hash, username, avatar_url, role, created_at, updated_at
		 FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *PostgresUserRepository) Update(user domain.User) (domain.User, error) {
	user.UpdatedAt = time.Now().UTC()
	tag, err := r.db.Exec(context.Background(),
		`UPDATE users SET email=$1, password_hash=$2, username=$3, avatar_url=$4, role=$5, updated_at=$6
		 WHERE id=$7`,
		user.Email, user.PasswordHash, user.Username, user.AvatarURL, string(user.Role), user.UpdatedAt, user.ID,
	)
	if err != nil {
		return domain.User{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func scanUser(row pgx.Row) (domain.User, bool, error) {
	var u domain.User
	var role string
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Username, &u.AvatarURL, &role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}
	u.Role = domain.Role(role)
	return u, true, nil
}
