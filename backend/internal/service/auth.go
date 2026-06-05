package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/maxwellpark/stanzabonanza/backend/internal/config"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

type AuthService struct {
	users      *repository.UserRepository
	sessions   *repository.SessionRepository
	magicLinks *repository.MagicLinkRepository
	cfg        *config.Config
}

func NewAuthService(users *repository.UserRepository, sessions *repository.SessionRepository, magicLinks *repository.MagicLinkRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		users:      users,
		sessions:   sessions,
		magicLinks: magicLinks,
		cfg:        cfg,
	}
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (uuid.UUID, error) {
	var hash = hashToken(token)
	session, err := s.sessions.GetByTokenHash(ctx, hash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid session: %w", err)
	}
	return session.UserID, nil
}

func (s *AuthService) RequestMagicLink(ctx context.Context, email string) (string, error) {
	_, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var user = &domain.User{
				Email:       email,
				DisplayName: email,
			}
			if err := s.users.Create(ctx, user); err != nil {
				return "", fmt.Errorf("creating user: %w", err)
			}
		} else {
			return "", fmt.Errorf("looking up user: %w", err)
		}
	}

	rawToken, err := generateToken(32)
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	var hash = hashToken(rawToken)
	var link = &domain.MagicLink{
		Email:     email,
		TokenHash: hash,
		ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
	}
	if err := s.magicLinks.Create(ctx, link); err != nil {
		return "", fmt.Errorf("storing magic link: %w", err)
	}

	return rawToken, nil
}

func (s *AuthService) VerifyMagicLink(ctx context.Context, token string) (string, *domain.User, error) {
	var hash = hashToken(token)
	link, err := s.magicLinks.GetByTokenHash(ctx, hash)
	if err != nil {
		return "", nil, fmt.Errorf("invalid magic link: %w", err)
	}

	if link.UsedAt != nil {
		return "", nil, fmt.Errorf("magic link already used")
	}
	if time.Now().UTC().After(link.ExpiresAt) {
		return "", nil, fmt.Errorf("magic link expired")
	}

	if err := s.magicLinks.MarkUsed(ctx, link.ID); err != nil {
		return "", nil, fmt.Errorf("marking link used: %w", err)
	}

	user, err := s.users.GetByEmail(ctx, link.Email)
	if err != nil {
		return "", nil, fmt.Errorf("finding user: %w", err)
	}

	sessionToken, err := s.CreateSession(ctx, user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("creating session: %w", err)
	}

	return sessionToken, user, nil
}

func (s *AuthService) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	rawToken, err := generateToken(32)
	if err != nil {
		return "", fmt.Errorf("generating session token: %w", err)
	}

	var hash = hashToken(rawToken)
	var session = &domain.Session{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().UTC().Add(30 * 24 * time.Hour),
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return "", fmt.Errorf("storing session: %w", err)
	}

	return rawToken, nil
}

func (s *AuthService) DeleteSession(ctx context.Context, token string) error {
	var hash = hashToken(token)
	session, err := s.sessions.GetByTokenHash(ctx, hash)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	return s.sessions.Delete(ctx, session.ID)
}

func (s *AuthService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, avatarURL string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.DisplayName = displayName
	user.Bio = bio
	user.AvatarURL = avatarURL

	return s.users.Update(ctx, user)
}

// TODO: WebAuthn — BeginRegistration, FinishRegistration, BeginLogin, FinishLogin

func hashToken(token string) string {
	var h = sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func generateToken(n int) (string, error) {
	var b = make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
