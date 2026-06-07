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
	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPoemService struct{ mock.Mock }

func (m *mockPoemService) Create(ctx context.Context, userID uuid.UUID, title, description string, format domain.PoemFormat, approvalMode domain.ApprovalMode, maxStanzas *int) (*domain.Poem, error) {
	args := m.Called(ctx, userID, title, description, format, approvalMode, maxStanzas)
	p, _ := args.Get(0).(*domain.Poem)
	return p, args.Error(1)
}
func (m *mockPoemService) Get(ctx context.Context, id uuid.UUID) (*domain.Poem, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Poem)
	return p, args.Error(1)
}
func (m *mockPoemService) List(ctx context.Context, page domain.PaginationParams, format, sort string) ([]domain.Poem, int, error) {
	args := m.Called(ctx, page, format, sort)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemService) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, userID, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemService) Update(ctx context.Context, userID, poemID uuid.UUID, title, description string) error {
	return m.Called(ctx, userID, poemID, title, description).Error(0)
}
func (m *mockPoemService) Delete(ctx context.Context, userID, poemID uuid.UUID) error {
	return m.Called(ctx, userID, poemID).Error(0)
}
func (m *mockPoemService) ListStanzas(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error) {
	args := m.Called(ctx, poemID)
	s, _ := args.Get(0).([]domain.Stanza)
	return s, args.Error(1)
}
func (m *mockPoemService) SubmitStanza(ctx context.Context, userID, poemID uuid.UUID, text, literaryDevice string) (*domain.Stanza, error) {
	args := m.Called(ctx, userID, poemID, text, literaryDevice)
	s, _ := args.Get(0).(*domain.Stanza)
	return s, args.Error(1)
}
func (m *mockPoemService) ReviewStanza(ctx context.Context, userID, poemID, stanzaID uuid.UUID, approved bool) error {
	return m.Called(ctx, userID, poemID, stanzaID, approved).Error(0)
}
func (m *mockPoemService) Feed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, userID, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemService) Explore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemService) HallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}

func newPoemHandler(svc poemService) *PoemHandler {
	return &PoemHandler{svc: svc}
}

// POST /poems tests

func TestPoemHandler_Create_ValidBody(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()
	poem := &domain.Poem{ID: uuid.New(), Title: "My Poem", Format: domain.FormatFreeVerse}
	svc.On("Create", mock.Anything, userID, "My Poem", "", domain.FormatFreeVerse, domain.ApprovalOpen, (*int)(nil)).Return(poem, nil)

	body := `{"title":"My Poem","format":"free_verse"}`
	r := httptest.NewRequest(http.MethodPost, "/poems", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp.Error)
}

func TestPoemHandler_Create_MissingTitle(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()
	body := `{"format":"free_verse"}`
	r := httptest.NewRequest(http.MethodPost, "/poems", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp.Error, "title")
}

func TestPoemHandler_Create_MissingFormat(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()
	body := `{"title":"My Poem"}`
	r := httptest.NewRequest(http.MethodPost, "/poems", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Create(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp.Error, "format")
}

// GET /poems/:poemID tests

func TestPoemHandler_Get_ValidID(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, Title: "Found"}
	svc.On("Get", mock.Anything, poemID).Return(poem, nil)

	r := httptest.NewRequest(http.MethodGet, "/poems/"+poemID.String(), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Get(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPoemHandler_Get_InvalidUUID(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	r := httptest.NewRequest(http.MethodGet, "/poems/bad-id", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", "bad-id")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Get(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// POST /poems/:poemID/stanzas tests

func TestPoemHandler_SubmitStanza_Valid(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()
	poemID := uuid.New()
	stanza := &domain.Stanza{ID: uuid.New(), Text: "A stanza"}
	svc.On("SubmitStanza", mock.Anything, userID, poemID, "A stanza", "").Return(stanza, nil)

	body := `{"text":"A stanza"}`
	r := httptest.NewRequest(http.MethodPost, "/poems/"+poemID.String()+"/stanzas", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.SubmitStanza(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPoemHandler_SubmitStanza_EmptyText(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()
	poemID := uuid.New()

	body := `{"text":""}`
	r := httptest.NewRequest(http.MethodPost, "/poems/"+poemID.String()+"/stanzas", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.SubmitStanza(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp.Error, "text")
}

// PUT /poems/:poemID/stanzas/:stanzaID tests

func TestPoemHandler_ReviewStanza_Approve(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()
	poemID := uuid.New()
	stanzaID := uuid.New()
	svc.On("ReviewStanza", mock.Anything, userID, poemID, stanzaID, true).Return(nil)

	body := `{"approved":true}`
	r := httptest.NewRequest(http.MethodPut, "/poems/"+poemID.String()+"/stanzas/"+stanzaID.String(), bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", poemID.String())
	rctx.URLParams.Add("stanzaID", stanzaID.String())
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.ReviewStanza(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPoemHandler_ReviewStanza_InvalidUUID(t *testing.T) {
	svc := &mockPoemService{}
	h := newPoemHandler(svc)

	userID := uuid.New()

	body := `{"approved":true}`
	r := httptest.NewRequest(http.MethodPut, "/poems/bad-id/stanzas/also-bad", bytes.NewBufferString(body))
	r.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("poemID", "bad-id")
	rctx.URLParams.Add("stanzaID", "also-bad")
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.ReviewStanza(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	_ = errors.New("")
}
