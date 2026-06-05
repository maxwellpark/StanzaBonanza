package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type MagicLinkRepository struct {
	pool *pgxpool.Pool
}

func NewMagicLinkRepository(pool *pgxpool.Pool) *MagicLinkRepository {
	return &MagicLinkRepository{pool: pool}
}

func (r *MagicLinkRepository) Create(ctx context.Context, link *domain.MagicLink) error {
	if link.ID == uuid.Nil {
		link.ID = uuid.New()
	}
	link.CreatedAt = time.Now().UTC()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO magic_links (id, email, token_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		link.ID, link.Email, link.TokenHash, link.ExpiresAt, link.CreatedAt,
	)
	return err
}

func (r *MagicLinkRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.MagicLink, error) {
	var ml domain.MagicLink
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, token_hash, expires_at, used_at, created_at
		 FROM magic_links WHERE token_hash = $1`, tokenHash,
	).Scan(&ml.ID, &ml.Email, &ml.TokenHash, &ml.ExpiresAt, &ml.UsedAt, &ml.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ml, nil
}

func (r *MagicLinkRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE magic_links SET used_at = $1 WHERE id = $2`,
		time.Now().UTC(), id,
	)
	return err
}
