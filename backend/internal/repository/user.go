package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	var now = time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, display_name, email, bio, avatar_url, is_verified, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.DisplayName, user.Email, user.Bio, user.AvatarURL, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, display_name, email, bio, avatar_url, is_verified, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.DisplayName, &u.Email, &u.Bio, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, display_name, email, bio, avatar_url, is_verified, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.DisplayName, &u.Email, &u.Bio, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET display_name = $1, bio = $2, avatar_url = $3, updated_at = $4
		 WHERE id = $5`,
		user.DisplayName, user.Bio, user.AvatarURL, user.UpdatedAt, user.ID,
	)
	return err
}
