package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSocialService struct{ mock.Mock }

func (m *mockSocialService) ToggleLike(ctx context.Context, userID, poemID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, poemID)
	return args.Bool(0), args.Error(1)
}
func (m *mockSocialService) AddComment(ctx context.Context, userID, poemID uuid.UUID, parentID *uuid.UUID, text string) (*domain.Comment, error) {
	args := m.Called(ctx, userID, poemID, parentID, text)
	c, _ := args.Get(0).(*domain.Comment)
	return c, args.Error(1)
}
func (m *mockSocialService) DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error {
	return m.Called(ctx, userID, commentID).Error(0)
}
func (m *mockSocialService) ListComments(ctx context.Context, poemID uuid.UUID, page domain.PaginationParams) ([]domain.Comment, int, error) {
	args := m.Called(ctx, poemID, page)
	c, _ := args.Get(0).([]domain.Comment)
	return c, args.Int(1), args.Error(2)
}
func (m *mockSocialService) ToggleFollow(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0), args.Error(1)
}
func (m *mockSocialService) ListFollowers(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	args := m.Called(ctx, userID, page)
	u, _ := args.Get(0).([]domain.User)
	return u, args.Int(1), args.Error(2)
}
func (m *mockSocialService) ListFollowing(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	args := m.Called(ctx, userID, page)
	u, _ := args.Get(0).([]domain.User)
	return u, args.Int(1), args.Error(2)
}
func (m *mockSocialService) ListNotifications(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Notification, int, error) {
	args := m.Called(ctx, userID, page)
	n, _ := args.Get(0).([]domain.Notification)
	return n, args.Int(1), args.Error(2)
}
func (m *mockSocialService) MarkNotificationsRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return m.Called(ctx, userID, ids).Error(0)
}

func newSocialHandler(svc socialService) *SocialHandler {
	return &SocialHandler{svc: svc}
}

// POST /poems/:poemID/like tests

func TestSocialHandler_ToggleLike_Valid(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	userID := uuid.New()
	poemID := uuid.New()
	svc.On("ToggleLike", mock.Anything, userID, poemID).Return(true, nil)

	r := httptest.NewRequest(http.MethodPost, "/poems/"+poemID.String()+"/like", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.ToggleLike(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp.Error)
	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, true, data["liked"])
}

// POST /poems/:poemID/comments tests

func TestSocialHandler_AddComment_Valid(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	userID := uuid.New()
	poemID := uuid.New()
	comment := &domain.Comment{ID: uuid.New(), Text: "Nice poem"}
	svc.On("AddComment", mock.Anything, userID, poemID, (*uuid.UUID)(nil), "Nice poem").Return(comment, nil)

	body := `{"text":"Nice poem"}`
	r := httptest.NewRequest(http.MethodPost, "/poems/"+poemID.String()+"/comments", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.AddComment(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSocialHandler_AddComment_EmptyText(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	userID := uuid.New()
	poemID := uuid.New()

	body := `{"text":""}`
	r := httptest.NewRequest(http.MethodPost, "/poems/"+poemID.String()+"/comments", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.AddComment(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Error)
}

// POST /users/:userID/follow tests

func TestSocialHandler_ToggleFollow_Valid(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	followerID := uuid.New()
	followedID := uuid.New()
	svc.On("ToggleFollow", mock.Anything, followerID, followedID).Return(true, nil)

	r := httptest.NewRequest(http.MethodPost, "/users/"+followedID.String()+"/follow", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", followedID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, followerID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.ToggleFollow(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, true, data["following"])
}

// GET /notifications tests

func TestSocialHandler_ListNotifications_Authenticated(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	userID := uuid.New()
	notifs := []domain.Notification{{ID: uuid.New(), RecipientID: userID}}
	svc.On("ListNotifications", mock.Anything, userID, mock.AnythingOfType("domain.PaginationParams")).Return(notifs, 1, nil)

	r := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.ListNotifications(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp.Error)
}

// POST /notifications/read tests

func TestSocialHandler_MarkNotificationsRead_ValidIDs(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	userID := uuid.New()
	id1 := uuid.New()
	ids := []uuid.UUID{id1}
	svc.On("MarkNotificationsRead", mock.Anything, userID, ids).Return(nil)

	bodyMap := map[string]interface{}{"ids": []string{id1.String()}}
	bodyBytes, _ := json.Marshal(bodyMap)
	r := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBuffer(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.MarkNotificationsRead(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSocialHandler_MarkNotificationsRead_EmptyIDs(t *testing.T) {
	svc := &mockSocialService{}
	h := newSocialHandler(svc)

	userID := uuid.New()

	body := `{"ids":[]}`
	r := httptest.NewRequest(http.MethodPost, "/notifications/read", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.MarkNotificationsRead(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
