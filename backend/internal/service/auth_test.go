package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockUserStore

type mockUserStore struct{ mock.Mock }

func (m *mockUserStore) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*domain.User)
	return u, args.Error(1)
}
func (m *mockUserStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*domain.User)
	return u, args.Error(1)
}
func (m *mockUserStore) Update(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

// mockSessionStore

type mockSessionStore struct{ mock.Mock }

func (m *mockSessionStore) Create(ctx context.Context, session *domain.Session) error {
	return m.Called(ctx, session).Error(0)
}
func (m *mockSessionStore) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	args := m.Called(ctx, tokenHash)
	s, _ := args.Get(0).(*domain.Session)
	return s, args.Error(1)
}
func (m *mockSessionStore) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// mockMagicLinkStore

type mockMagicLinkStore struct{ mock.Mock }

func (m *mockMagicLinkStore) Create(ctx context.Context, link *domain.MagicLink) error {
	return m.Called(ctx, link).Error(0)
}
func (m *mockMagicLinkStore) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.MagicLink, error) {
	args := m.Called(ctx, tokenHash)
	l, _ := args.Get(0).(*domain.MagicLink)
	return l, args.Error(1)
}
func (m *mockMagicLinkStore) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

// mockWebAuthnStore

type mockWebAuthnStore struct{ mock.Mock }

func (m *mockWebAuthnStore) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.WebAuthnCredential, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]domain.WebAuthnCredential)
	return c, args.Error(1)
}
func (m *mockWebAuthnStore) Create(ctx context.Context, cred *domain.WebAuthnCredential) error {
	return m.Called(ctx, cred).Error(0)
}
func (m *mockWebAuthnStore) GetByCredentialID(ctx context.Context, credentialID []byte) (*domain.WebAuthnCredential, error) {
	args := m.Called(ctx, credentialID)
	c, _ := args.Get(0).(*domain.WebAuthnCredential)
	return c, args.Error(1)
}
func (m *mockWebAuthnStore) UpdateSignCount(ctx context.Context, credentialID []byte, signCount uint32) error {
	return m.Called(ctx, credentialID, signCount).Error(0)
}

func newAuthSvc(users *mockUserStore, sessions *mockSessionStore, links *mockMagicLinkStore) *AuthService {
	return &AuthService{
		users:      users,
		sessions:   sessions,
		magicLinks: links,
		webAuthns:  &mockWebAuthnStore{},
		waSessions: make(map[string]*waSession),
	}
}

// RequestMagicLink tests

func TestAuthService_RequestMagicLink_NewUser(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	users.On("GetByEmail", mock.Anything, "new@test.com").Return(nil, pgx.ErrNoRows)
	users.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	links.On("Create", mock.Anything, mock.AnythingOfType("*domain.MagicLink")).Return(nil)

	token, err := svc.RequestMagicLink(context.Background(), "new@test.com")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	users.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*domain.User"))
}

func TestAuthService_RequestMagicLink_ExistingUser(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	existing := &domain.User{ID: uuid.New(), Email: "old@test.com"}
	users.On("GetByEmail", mock.Anything, "old@test.com").Return(existing, nil)
	links.On("Create", mock.Anything, mock.AnythingOfType("*domain.MagicLink")).Return(nil)

	token, err := svc.RequestMagicLink(context.Background(), "old@test.com")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	users.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// VerifyMagicLink tests

func TestAuthService_VerifyMagicLink_Success(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	rawToken, _ := generateToken(32)
	hash := hashToken(rawToken)
	linkID := uuid.New()
	userID := uuid.New()

	link := &domain.MagicLink{
		ID:        linkID,
		Email:     "user@test.com",
		TokenHash: hash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	user := &domain.User{ID: userID, Email: "user@test.com"}

	links.On("GetByTokenHash", mock.Anything, hash).Return(link, nil)
	links.On("MarkUsed", mock.Anything, linkID).Return(nil)
	users.On("GetByEmail", mock.Anything, "user@test.com").Return(user, nil)
	sessions.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(nil)

	sessionToken, returnedUser, err := svc.VerifyMagicLink(context.Background(), rawToken)
	require.NoError(t, err)
	assert.NotEmpty(t, sessionToken)
	assert.Equal(t, userID, returnedUser.ID)
}

func TestAuthService_VerifyMagicLink_AlreadyUsed(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	rawToken, _ := generateToken(32)
	hash := hashToken(rawToken)
	usedAt := time.Now().Add(-1 * time.Minute)

	link := &domain.MagicLink{
		ID:        uuid.New(),
		TokenHash: hash,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		UsedAt:    &usedAt,
	}
	links.On("GetByTokenHash", mock.Anything, hash).Return(link, nil)

	_, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already used")
}

func TestAuthService_VerifyMagicLink_Expired(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	rawToken, _ := generateToken(32)
	hash := hashToken(rawToken)

	link := &domain.MagicLink{
		ID:        uuid.New(),
		TokenHash: hash,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	links.On("GetByTokenHash", mock.Anything, hash).Return(link, nil)

	_, _, err := svc.VerifyMagicLink(context.Background(), rawToken)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestAuthService_VerifyMagicLink_InvalidToken(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	links.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	_, _, err := svc.VerifyMagicLink(context.Background(), "badtoken")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid magic link")
}

// CreateSession tests

func TestAuthService_CreateSession_Success(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	userID := uuid.New()
	sessions.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(nil)

	token, err := svc.CreateSession(context.Background(), userID)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// DeleteSession tests

func TestAuthService_DeleteSession_Success(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	sessionID := uuid.New()
	rawToken, _ := generateToken(32)
	hash := hashToken(rawToken)

	sessions.On("GetByTokenHash", mock.Anything, hash).Return(&domain.Session{ID: sessionID}, nil)
	sessions.On("Delete", mock.Anything, sessionID).Return(nil)

	err := svc.DeleteSession(context.Background(), rawToken)
	require.NoError(t, err)
}

func TestAuthService_DeleteSession_NotFound(t *testing.T) {
	users := &mockUserStore{}
	sessions := &mockSessionStore{}
	links := &mockMagicLinkStore{}
	svc := newAuthSvc(users, sessions, links)

	sessions.On("GetByTokenHash", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	err := svc.DeleteSession(context.Background(), "notoken")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

// UpdateProfile tests

func TestAuthService_UpdateProfile_Success(t *testing.T) {
	users := &mockUserStore{}
	svc := newAuthSvc(users, &mockSessionStore{}, &mockMagicLinkStore{})

	userID := uuid.New()
	user := &domain.User{ID: userID, DisplayName: "old"}

	users.On("GetByID", mock.Anything, userID).Return(user, nil)
	users.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.UpdateProfile(context.Background(), userID, "new name", "bio", "http://avatar.png")
	require.NoError(t, err)
	assert.Equal(t, "new name", user.DisplayName)
}

func TestAuthService_UpdateProfile_UserNotFound(t *testing.T) {
	users := &mockUserStore{}
	svc := newAuthSvc(users, &mockSessionStore{}, &mockMagicLinkStore{})

	users.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	err := svc.UpdateProfile(context.Background(), uuid.New(), "name", "", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

// hashToken / generateToken tests

func TestHashToken_Deterministic(t *testing.T) {
	h1 := hashToken("mysecret")
	h2 := hashToken("mysecret")
	assert.Equal(t, h1, h2)
	assert.NotEmpty(t, h1)
}

func TestGenerateToken_UniqueOutputs(t *testing.T) {
	t1, err1 := generateToken(32)
	t2, err2 := generateToken(32)
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, t1, t2)
	// 32 bytes hex-encoded = 64 chars.
	assert.Len(t, t1, 64)
}
