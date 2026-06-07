package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type TutorialRepository struct {
	pool *pgxpool.Pool
}

func NewTutorialRepository(pool *pgxpool.Pool) *TutorialRepository {
	return &TutorialRepository{pool: pool}
}

func (r *TutorialRepository) List(ctx context.Context) ([]domain.Tutorial, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, slug, format, content_md, difficulty, display_order, created_at
		 FROM tutorials ORDER BY display_order ASC, created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tutorials []domain.Tutorial
	for rows.Next() {
		var t domain.Tutorial
		if err := rows.Scan(&t.ID, &t.Title, &t.Slug, &t.Format, &t.ContentMD, &t.Difficulty, &t.DisplayOrder, &t.CreatedAt); err != nil {
			return nil, err
		}
		tutorials = append(tutorials, t)
	}
	return tutorials, rows.Err()
}

func (r *TutorialRepository) GetBySlug(ctx context.Context, slug string) (*domain.Tutorial, error) {
	var t domain.Tutorial
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, slug, format, content_md, difficulty, display_order, created_at
		 FROM tutorials WHERE slug = $1`, slug,
	).Scan(&t.ID, &t.Title, &t.Slug, &t.Format, &t.ContentMD, &t.Difficulty, &t.DisplayOrder, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
