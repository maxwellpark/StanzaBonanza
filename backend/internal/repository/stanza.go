package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type StanzaRepository struct {
	pool *pgxpool.Pool
}

func NewStanzaRepository(pool *pgxpool.Pool) *StanzaRepository {
	return &StanzaRepository{pool: pool}
}

func (r *StanzaRepository) Create(ctx context.Context, stanza *domain.Stanza) error {
	if stanza.ID == uuid.Nil {
		stanza.ID = uuid.New()
	}
	stanza.CreatedAt = time.Now().UTC()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO stanzas (id, poem_id, author_id, text, position, literary_device, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		stanza.ID, stanza.PoemID, stanza.AuthorID, stanza.Text, stanza.Position,
		stanza.LiteraryDevice, stanza.Status, stanza.CreatedAt,
	)
	return err
}

func (r *StanzaRepository) ListByPoem(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.poem_id, s.author_id, s.text, s.position, s.literary_device, s.status, s.created_at,
		        u.id, u.display_name, u.email, u.bio, u.avatar_url, u.is_verified, u.created_at, u.updated_at
		 FROM stanzas s JOIN users u ON u.id = s.author_id
		 WHERE s.poem_id = $1 ORDER BY s.position ASC`, poemID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stanzas := make([]domain.Stanza, 0)
	for rows.Next() {
		var s domain.Stanza
		var author domain.User
		err := rows.Scan(
			&s.ID, &s.PoemID, &s.AuthorID, &s.Text, &s.Position, &s.LiteraryDevice, &s.Status, &s.CreatedAt,
			&author.ID, &author.DisplayName, &author.Email, &author.Bio, &author.AvatarURL,
			&author.IsVerified, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		s.Author = &author
		stanzas = append(stanzas, s)
	}
	return stanzas, rows.Err()
}

func (r *StanzaRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.StanzaStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE stanzas SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *StanzaRepository) GetNextPosition(ctx context.Context, poemID uuid.UUID) (int, error) {
	var maxPos *int
	err := r.pool.QueryRow(ctx, `SELECT MAX(position) FROM stanzas WHERE poem_id = $1`, poemID).Scan(&maxPos)
	if err != nil {
		return 0, err
	}
	if maxPos == nil {
		return 1, nil
	}
	return *maxPos + 1, nil
}
