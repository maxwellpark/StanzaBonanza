package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
)

type WebAuthnRepository struct {
	pool *pgxpool.Pool
}

func NewWebAuthnRepository(pool *pgxpool.Pool) *WebAuthnRepository {
	return &WebAuthnRepository{pool: pool}
}

func (r *WebAuthnRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.WebAuthnCredential, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, credential_id, public_key, sign_count, transports, created_at
		 FROM webauthn_credentials WHERE user_id = $1`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []domain.WebAuthnCredential
	for rows.Next() {
		var c domain.WebAuthnCredential
		if err := rows.Scan(&c.ID, &c.UserID, &c.CredentialID, &c.PublicKey, &c.SignCount, &c.Transports, &c.CreatedAt); err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}
	return creds, rows.Err()
}

func (r *WebAuthnRepository) Create(ctx context.Context, cred *domain.WebAuthnCredential) error {
	if cred.ID == uuid.Nil {
		cred.ID = uuid.New()
	}
	cred.CreatedAt = time.Now().UTC()

	_, err := r.pool.Exec(ctx,
		`INSERT INTO webauthn_credentials (id, user_id, credential_id, public_key, sign_count, transports, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		cred.ID, cred.UserID, cred.CredentialID, cred.PublicKey, cred.SignCount, cred.Transports, cred.CreatedAt,
	)
	return err
}

func (r *WebAuthnRepository) GetByCredentialID(ctx context.Context, credentialID []byte) (*domain.WebAuthnCredential, error) {
	var c domain.WebAuthnCredential
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, credential_id, public_key, sign_count, transports, created_at
		 FROM webauthn_credentials WHERE credential_id = $1`, credentialID,
	).Scan(&c.ID, &c.UserID, &c.CredentialID, &c.PublicKey, &c.SignCount, &c.Transports, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *WebAuthnRepository) UpdateSignCount(ctx context.Context, credentialID []byte, signCount uint32) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE webauthn_credentials SET sign_count = $1 WHERE credential_id = $2`,
		signCount, credentialID,
	)
	return err
}
