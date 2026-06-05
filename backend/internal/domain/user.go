package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatarUrl"`
	IsVerified  bool      `json:"isVerified"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type WebAuthnCredential struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	CredentialID []byte    `json:"-"`
	PublicKey    []byte    `json:"-"`
	SignCount    uint32    `json:"-"`
	Transports   []string  `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type MagicLink struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expiresAt"`
	UsedAt    *time.Time `json:"usedAt"`
	CreatedAt time.Time  `json:"createdAt"`
}

type UserTier struct {
	UserID    uuid.UUID  `json:"userId"`
	Tier      string     `json:"tier"`
	ExpiresAt *time.Time `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
}
