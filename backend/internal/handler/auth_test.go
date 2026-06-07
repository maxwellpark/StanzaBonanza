package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAuthService struct{ mock.Mock }

func (m *mockAuthService) RequestMagicLink(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}
func (m *mockAuthService) VerifyMagicLink(ctx context.Context, token string) (string, *domain.User, error) {
	args := m.Called(ctx, token)
	u, _ := args.Get(1).(*domain.User)
	return args.String(0), u, args.Error(2)
}
func (m *mockAuthService) DeleteSession(ctx context.Context, token string) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockAuthService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	u, _ := args.Get(0).(*domain.User)
	return u, args.Error(1)
}
func (m *mockAuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, avatarURL string) error {
	return m.Called(ctx, userID, displayName, bio, avatarURL).Error(0)
}
func (m *mockAuthService) BeginRegistration(ctx context.Context, userID uuid.UUID) (*protocol.CredentialCreation, string, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).(*protocol.CredentialCreation)
	return c, args.String(1), args.Error(2)
}
func (m *mockAuthService) FinishRegistration(ctx context.Context, userID uuid.UUID, sessionKey string, r *http.Request) (*domain.WebAuthnCredential, error) {
	args := m.Called(ctx, userID, sessionKey, r)
	c, _ := args.Get(0).(*domain.WebAuthnCredential)
	return c, args.Error(1)
}
func (m *mockAuthService) BeginLogin(ctx context.Context) (*protocol.CredentialAssertion, string, error) {
	args := m.Called(ctx)
	a, _ := args.Get(0).(*protocol.CredentialAssertion)
	return a, args.String(1), args.Error(2)
}
func (m *mockAuthService) FinishLogin(ctx context.Context, sessionKey string, r *http.Request) (*domain.User, string, error) {
	args := m.Called(ctx, sessionKey, r)
	u, _ := args.Get(0).(*domain.User)
	return u, args.String(1), args.Error(2)
}

func newAuthHandler(svc authService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// POST /auth/magic-link tests

func TestAuthHandler_RequestMagicLink_ValidEmail(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	svc.On("RequestMagicLink", mock.Anything, "user@test.com").Return("rawtoken", nil)

	body := `{"email":"user@test.com"}`
	r := httptest.NewRequest(http.MethodPost, "/auth/magic-link", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RequestMagicLink(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp.Error)
}

func TestAuthHandler_RequestMagicLink_MissingEmail(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	body := `{"email":""}`
	r := httptest.NewRequest(http.MethodPost, "/auth/magic-link", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RequestMagicLink(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Error)
}

// GET /auth/magic-link/verify tests

func TestAuthHandler_VerifyMagicLink_ValidToken(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "user@test.com"}
	svc.On("VerifyMagicLink", mock.Anything, "validtoken").Return("sessiontoken", user, nil)

	r := httptest.NewRequest(http.MethodGet, "/auth/magic-link/verify?token=validtoken", nil)
	w := httptest.NewRecorder()

	h.VerifyMagicLink(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	// Session cookie must be set.
	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "session" {
			found = true
			assert.Equal(t, "sessiontoken", c.Value)
		}
	}
	assert.True(t, found, "expected session cookie to be set")
}

func TestAuthHandler_VerifyMagicLink_MissingToken(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/auth/magic-link/verify", nil)
	w := httptest.NewRecorder()

	h.VerifyMagicLink(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// POST /auth/logout tests

func TestAuthHandler_Logout_ClearsCookie(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	svc.On("DeleteSession", mock.Anything, "mysession").Return(nil)

	r := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "mysession"})
	w := httptest.NewRecorder()

	h.Logout(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "session" {
			found = true
			assert.Equal(t, -1, c.MaxAge)
		}
	}
	assert.True(t, found, "expected session cookie to be cleared")
}

// GET /auth/me tests

func TestAuthHandler_Me_Authenticated(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "me@test.com"}
	svc.On("GetUser", mock.Anything, userID).Return(user, nil)

	r := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Me(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp.Error)
}

func TestAuthHandler_Me_Unauthenticated(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	w := httptest.NewRecorder()

	h.Me(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// GET /users/:userID tests

func TestAuthHandler_GetProfile_ValidUUID(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "profile@test.com"}
	svc.On("GetUser", mock.Anything, userID).Return(user, nil)

	r := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", userID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.GetProfile(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetProfile_InvalidUUID(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/users/not-a-uuid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "not-a-uuid")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.GetProfile(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// PUT /auth/me tests

func TestAuthHandler_UpdateProfile_Success(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	userID := uuid.New()
	svc.On("UpdateProfile", mock.Anything, userID, "New Name", "bio text", "http://avatar.png").Return(nil)

	body := `{"displayName":"New Name","bio":"bio text","avatarUrl":"http://avatar.png"}`
	r := httptest.NewRequest(http.MethodPut, "/auth/me", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateProfile(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_UpdateProfile_DisplayNameTooLong(t *testing.T) {
	svc := &mockAuthService{}
	h := newAuthHandler(svc)

	userID := uuid.New()
	longName := string(make([]byte, 51))
	body, _ := json.Marshal(map[string]string{"displayName": longName})
	r := httptest.NewRequest(http.MethodPut, "/auth/me", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.UpdateProfile(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp.Error, "50")

	// suppress "declared but not used" if svc is only setup - force call check
	svc.AssertNotCalled(t, "UpdateProfile")
	_ = errors.New("")
}
