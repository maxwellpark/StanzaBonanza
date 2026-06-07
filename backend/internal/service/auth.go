package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/maxwellpark/stanzabonanza/backend/internal/config"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

// webAuthnUser wraps domain.User to satisfy the webauthn.User interface.
type webAuthnUser struct {
	user  *domain.User
	creds []webauthn.Credential
}

func (u *webAuthnUser) WebAuthnID() []byte            { return u.user.ID[:] }
func (u *webAuthnUser) WebAuthnName() string           { return u.user.Email }
func (u *webAuthnUser) WebAuthnDisplayName() string    { return u.user.DisplayName }
func (u *webAuthnUser) WebAuthnCredentials() []webauthn.Credential { return u.creds }

// waSession holds a pending WebAuthn challenge with an expiry.
type waSession struct {
	data      *webauthn.SessionData
	expiresAt time.Time
}

type userStore interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type sessionStore interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type magicLinkStore interface {
	Create(ctx context.Context, link *domain.MagicLink) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.MagicLink, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
}

type webAuthnStore interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.WebAuthnCredential, error)
	Create(ctx context.Context, cred *domain.WebAuthnCredential) error
	GetByCredentialID(ctx context.Context, credentialID []byte) (*domain.WebAuthnCredential, error)
	UpdateSignCount(ctx context.Context, credentialID []byte, signCount uint32) error
}

type AuthService struct {
	users      userStore
	sessions   sessionStore
	magicLinks magicLinkStore
	webAuthns  webAuthnStore
	wac        *webauthn.WebAuthn
	cfg        *config.Config

	waMu       sync.Mutex
	waSessions map[string]*waSession
}

func NewAuthService(
	users *repository.UserRepository,
	sessions *repository.SessionRepository,
	magicLinks *repository.MagicLinkRepository,
	webAuthns *repository.WebAuthnRepository,
	cfg *config.Config,
) *AuthService {
	wac, _ := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.WebAuthnRPName,
		RPID:          cfg.WebAuthnRPID,
		RPOrigins:     cfg.WebAuthnOrigins,
	})
	return &AuthService{
		users:      users,
		sessions:   sessions,
		magicLinks: magicLinks,
		webAuthns:  webAuthns,
		wac:        wac,
		cfg:        cfg,
		waSessions: make(map[string]*waSession),
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

func (s *AuthService) storeWASession(key string, data *webauthn.SessionData) {
	s.waMu.Lock()
	defer s.waMu.Unlock()
	s.waSessions[key] = &waSession{data: data, expiresAt: time.Now().Add(5 * time.Minute)}
}

func (s *AuthService) popWASession(key string) (*webauthn.SessionData, bool) {
	s.waMu.Lock()
	defer s.waMu.Unlock()
	sess, ok := s.waSessions[key]
	if !ok || time.Now().After(sess.expiresAt) {
		delete(s.waSessions, key)
		return nil, false
	}
	delete(s.waSessions, key)
	return sess.data, true
}

func (s *AuthService) buildWebAuthnUser(ctx context.Context, user *domain.User) (*webAuthnUser, error) {
	domainCreds, err := s.webAuthns.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	creds := make([]webauthn.Credential, 0, len(domainCreds))
	for _, c := range domainCreds {
		creds = append(creds, webauthn.Credential{
			ID:        c.CredentialID,
			PublicKey: c.PublicKey,
			Authenticator: webauthn.Authenticator{
				SignCount: c.SignCount,
			},
		})
	}
	return &webAuthnUser{user: user, creds: creds}, nil
}

func (s *AuthService) BeginRegistration(ctx context.Context, userID uuid.UUID) (*protocol.CredentialCreation, string, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	waUser, err := s.buildWebAuthnUser(ctx, user)
	if err != nil {
		return nil, "", fmt.Errorf("building webauthn user: %w", err)
	}

	creation, sessionData, err := s.wac.BeginRegistration(waUser)
	if err != nil {
		return nil, "", fmt.Errorf("begin registration: %w", err)
	}

	key, err := generateToken(16)
	if err != nil {
		return nil, "", fmt.Errorf("generating session key: %w", err)
	}
	s.storeWASession(key, sessionData)

	return creation, key, nil
}

func (s *AuthService) FinishRegistration(ctx context.Context, userID uuid.UUID, sessionKey string, r *http.Request) (*domain.WebAuthnCredential, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	waUser, err := s.buildWebAuthnUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("building webauthn user: %w", err)
	}

	sessionData, ok := s.popWASession(sessionKey)
	if !ok {
		return nil, fmt.Errorf("session expired or not found")
	}

	credential, err := s.wac.FinishRegistration(waUser, *sessionData, r)
	if err != nil {
		return nil, fmt.Errorf("finish registration: %w", err)
	}

	transports := make([]string, len(credential.Transport))
	for i, t := range credential.Transport {
		transports[i] = string(t)
	}
	cred := &domain.WebAuthnCredential{
		UserID:       userID,
		CredentialID: credential.ID,
		PublicKey:    credential.PublicKey,
		SignCount:    credential.Authenticator.SignCount,
		Transports:   transports,
	}
	if err := s.webAuthns.Create(ctx, cred); err != nil {
		return nil, fmt.Errorf("storing credential: %w", err)
	}

	return cred, nil
}

func (s *AuthService) BeginLogin(ctx context.Context) (*protocol.CredentialAssertion, string, error) {
	assertion, sessionData, err := s.wac.BeginDiscoverableLogin()
	if err != nil {
		return nil, "", fmt.Errorf("begin login: %w", err)
	}

	key, err := generateToken(16)
	if err != nil {
		return nil, "", fmt.Errorf("generating session key: %w", err)
	}
	s.storeWASession(key, sessionData)

	return assertion, key, nil
}

func (s *AuthService) FinishLogin(ctx context.Context, sessionKey string, r *http.Request) (*domain.User, string, error) {
	sessionData, ok := s.popWASession(sessionKey)
	if !ok {
		return nil, "", fmt.Errorf("session expired or not found")
	}

	userHandler := func(rawID, userHandle []byte) (webauthn.User, error) {
		userID, err := uuid.FromBytes(userHandle)
		if err != nil {
			return nil, fmt.Errorf("invalid user handle: %w", err)
		}
		user, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return s.buildWebAuthnUser(ctx, user)
	}

	credential, err := s.wac.FinishDiscoverableLogin(userHandler, *sessionData, r)
	if err != nil {
		return nil, "", fmt.Errorf("finish login: %w", err)
	}

	domainCred, err := s.webAuthns.GetByCredentialID(ctx, credential.ID)
	if err != nil {
		return nil, "", fmt.Errorf("credential not found: %w", err)
	}

	if err := s.webAuthns.UpdateSignCount(ctx, credential.ID, credential.Authenticator.SignCount); err != nil {
		return nil, "", fmt.Errorf("updating sign count: %w", err)
	}

	user, err := s.users.GetByID(ctx, domainCred.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	sessionToken, err := s.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("creating session: %w", err)
	}

	return user, sessionToken, nil
}

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
