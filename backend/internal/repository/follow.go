package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type FollowRepository struct {
	pool *pgxpool.Pool
}

func NewFollowRepository(pool *pgxpool.Pool) *FollowRepository {
	return &FollowRepository{pool: pool}
}

func (r *FollowRepository) Create(ctx context.Context, follow *domain.Follow) error {
	follow.CreatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`INSERT INTO follows (follower_id, followed_id, created_at) VALUES ($1, $2, $3)`,
		follow.FollowerID, follow.FollowedID, follow.CreatedAt,
	)
	return err
}

func (r *FollowRepository) Delete(ctx context.Context, followerID, followedID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM follows WHERE follower_id = $1 AND followed_id = $2`,
		followerID, followedID,
	)
	return err
}

func (r *FollowRepository) Exists(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND followed_id = $2)`,
		followerID, followedID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *FollowRepository) ListFollowers(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	page.Normalize()

	var totalCount int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM follows WHERE followed_id = $1`, userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.display_name, u.email, u.bio, u.avatar_url, u.is_verified, u.created_at, u.updated_at
		 FROM follows f
		 JOIN users u ON u.id = f.follower_id
		 WHERE f.followed_id = $1
		 ORDER BY f.created_at DESC LIMIT $2 OFFSET $3`,
		userID, page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.DisplayName, &u.Email, &u.Bio, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, totalCount, rows.Err()
}

func (r *FollowRepository) ListFollowing(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	page.Normalize()

	var totalCount int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM follows WHERE follower_id = $1`, userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.display_name, u.email, u.bio, u.avatar_url, u.is_verified, u.created_at, u.updated_at
		 FROM follows f
		 JOIN users u ON u.id = f.followed_id
		 WHERE f.follower_id = $1
		 ORDER BY f.created_at DESC LIMIT $2 OFFSET $3`,
		userID, page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.DisplayName, &u.Email, &u.Bio, &u.AvatarURL, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, totalCount, rows.Err()
}
