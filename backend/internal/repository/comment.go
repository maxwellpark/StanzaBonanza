package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	var now = time.Now().UTC()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	_, err := r.pool.Exec(ctx,
		`INSERT INTO comments (id, poem_id, author_id, parent_id, text, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		comment.ID, comment.PoemID, comment.AuthorID, comment.ParentID,
		comment.Text, comment.CreatedAt, comment.UpdatedAt,
	)
	return err
}

func (r *CommentRepository) ListByPoem(ctx context.Context, poemID uuid.UUID, page domain.PaginationParams) ([]domain.Comment, int, error) {
	page.Normalize()

	var totalCount int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM comments WHERE poem_id = $1`, poemID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT c.id, c.poem_id, c.author_id, c.parent_id, c.text, c.created_at, c.updated_at,
		        u.id, u.display_name, u.avatar_url
		 FROM comments c
		 JOIN users u ON u.id = c.author_id
		 WHERE c.poem_id = $1
		 ORDER BY c.created_at ASC LIMIT $2 OFFSET $3`,
		poemID, page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	comments := make([]domain.Comment, 0)
	for rows.Next() {
		var c domain.Comment
		var author domain.User
		if err := rows.Scan(
			&c.ID, &c.PoemID, &c.AuthorID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt,
			&author.ID, &author.DisplayName, &author.AvatarURL,
		); err != nil {
			return nil, 0, err
		}
		c.Author = &author
		comments = append(comments, c)
	}
	return comments, totalCount, rows.Err()
}

func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1`, id)
	return err
}

func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	var c domain.Comment
	err := r.pool.QueryRow(ctx,
		`SELECT id, poem_id, author_id, parent_id, text, created_at, updated_at
		 FROM comments WHERE id = $1`, id,
	).Scan(&c.ID, &c.PoemID, &c.AuthorID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
