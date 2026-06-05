package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type PoemRepository struct {
	pool *pgxpool.Pool
}

func NewPoemRepository(pool *pgxpool.Pool) *PoemRepository {
	return &PoemRepository{pool: pool}
}

func (r *PoemRepository) Create(ctx context.Context, poem *domain.Poem) error {
	if poem.ID == uuid.Nil {
		poem.ID = uuid.New()
	}
	var now = time.Now().UTC()
	poem.CreatedAt = now
	poem.UpdatedAt = now

	_, err := r.pool.Exec(ctx,
		`INSERT INTO poems (id, author_id, title, description, format, format_rules_json, approval_mode, max_stanzas, is_hall_of_fame, like_count, stanza_count, comment_count, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		poem.ID, poem.AuthorID, poem.Title, poem.Description, poem.Format, poem.FormatRulesJSON,
		poem.ApprovalMode, poem.MaxStanzas, poem.IsHallOfFame, poem.LikeCount, poem.StanzaCount,
		poem.CommentCount, poem.CreatedAt, poem.UpdatedAt,
	)
	return err
}

const poemSelectFields = `
	p.id, p.author_id, p.title, p.description, p.format, p.format_rules_json,
	p.approval_mode, p.max_stanzas, p.is_hall_of_fame, p.like_count, p.stanza_count,
	p.comment_count, p.created_at, p.updated_at,
	u.id, u.display_name, u.email, u.bio, u.avatar_url, u.is_verified, u.created_at, u.updated_at`

const poemFromJoin = ` FROM poems p JOIN users u ON u.id = p.author_id`

func scanPoem(row pgx.Row) (*domain.Poem, error) {
	var p domain.Poem
	var author domain.User
	err := row.Scan(
		&p.ID, &p.AuthorID, &p.Title, &p.Description, &p.Format, &p.FormatRulesJSON,
		&p.ApprovalMode, &p.MaxStanzas, &p.IsHallOfFame, &p.LikeCount, &p.StanzaCount,
		&p.CommentCount, &p.CreatedAt, &p.UpdatedAt,
		&author.ID, &author.DisplayName, &author.Email, &author.Bio, &author.AvatarURL,
		&author.IsVerified, &author.CreatedAt, &author.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.Author = &author
	return &p, nil
}

func scanPoems(rows pgx.Rows) ([]domain.Poem, error) {
	var poems []domain.Poem
	for rows.Next() {
		var p domain.Poem
		var author domain.User
		err := rows.Scan(
			&p.ID, &p.AuthorID, &p.Title, &p.Description, &p.Format, &p.FormatRulesJSON,
			&p.ApprovalMode, &p.MaxStanzas, &p.IsHallOfFame, &p.LikeCount, &p.StanzaCount,
			&p.CommentCount, &p.CreatedAt, &p.UpdatedAt,
			&author.ID, &author.DisplayName, &author.Email, &author.Bio, &author.AvatarURL,
			&author.IsVerified, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.Author = &author
		poems = append(poems, p)
	}
	return poems, rows.Err()
}

func (r *PoemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Poem, error) {
	var query = `SELECT` + poemSelectFields + poemFromJoin + ` WHERE p.id = $1`
	return scanPoem(r.pool.QueryRow(ctx, query, id))
}

func (r *PoemRepository) List(ctx context.Context, page domain.PaginationParams, format string, sort string) ([]domain.Poem, int, error) {
	page.Normalize()

	var where string
	var args []any
	var argIdx = 1

	if format != "" {
		where = fmt.Sprintf(" WHERE p.format = $%d", argIdx)
		args = append(args, format)
		argIdx++
	}

	var orderBy string
	switch sort {
	case "popular":
		orderBy = " ORDER BY p.like_count DESC, p.created_at DESC"
	default:
		orderBy = " ORDER BY p.created_at DESC"
	}

	var countQuery = `SELECT COUNT(*)` + poemFromJoin + where
	var totalCount int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	var dataQuery = `SELECT` + poemSelectFields + poemFromJoin + where + orderBy +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, page.PageSize, page.Offset())

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	poems, err := scanPoems(rows)
	if err != nil {
		return nil, 0, err
	}
	return poems, totalCount, nil
}

func (r *PoemRepository) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	page.Normalize()

	var totalCount int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*)`+poemFromJoin+` WHERE p.author_id = $1`, userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT`+poemSelectFields+poemFromJoin+` WHERE p.author_id = $1 ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`,
		userID, page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	poems, err := scanPoems(rows)
	if err != nil {
		return nil, 0, err
	}
	return poems, totalCount, nil
}

func (r *PoemRepository) ListFeed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	page.Normalize()

	var feedWhere = poemFromJoin + ` WHERE p.author_id IN (SELECT followed_id FROM follows WHERE follower_id = $1)`

	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)`+feedWhere, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT`+poemSelectFields+feedWhere+` ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`,
		userID, page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	poems, err := scanPoems(rows)
	if err != nil {
		return nil, 0, err
	}
	return poems, totalCount, nil
}

func (r *PoemRepository) ListExplore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	page.Normalize()

	var exploreWhere = poemFromJoin + ` WHERE p.created_at > now() - interval '7 days'`

	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)`+exploreWhere).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT`+poemSelectFields+exploreWhere+` ORDER BY p.like_count DESC, p.created_at DESC LIMIT $1 OFFSET $2`,
		page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	poems, err := scanPoems(rows)
	if err != nil {
		return nil, 0, err
	}
	return poems, totalCount, nil
}

func (r *PoemRepository) ListHallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	page.Normalize()

	var hofWhere = poemFromJoin + ` WHERE p.is_hall_of_fame = true`

	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*)`+hofWhere).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT`+poemSelectFields+hofWhere+` ORDER BY p.like_count DESC, p.created_at DESC LIMIT $1 OFFSET $2`,
		page.PageSize, page.Offset(),
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	poems, err := scanPoems(rows)
	if err != nil {
		return nil, 0, err
	}
	return poems, totalCount, nil
}

func (r *PoemRepository) Update(ctx context.Context, poem *domain.Poem) error {
	poem.UpdatedAt = time.Now().UTC()
	_, err := r.pool.Exec(ctx,
		`UPDATE poems SET title = $1, description = $2, format = $3, format_rules_json = $4,
		 approval_mode = $5, max_stanzas = $6, is_hall_of_fame = $7, updated_at = $8
		 WHERE id = $9`,
		poem.Title, poem.Description, poem.Format, poem.FormatRulesJSON,
		poem.ApprovalMode, poem.MaxStanzas, poem.IsHallOfFame, poem.UpdatedAt, poem.ID,
	)
	return err
}

func (r *PoemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM poems WHERE id = $1`, id)
	return err
}

func (r *PoemRepository) IncrementCounter(ctx context.Context, id uuid.UUID, column string, delta int) error {
	var allowed = map[string]bool{
		"like_count":    true,
		"stanza_count":  true,
		"comment_count": true,
	}
	if !allowed[column] {
		return fmt.Errorf("invalid counter column: %s", column)
	}
	var query = fmt.Sprintf(`UPDATE poems SET %s = %s + $1, updated_at = now() WHERE id = $2`, column, column)
	_, err := r.pool.Exec(ctx, query, delta, id)
	return err
}
