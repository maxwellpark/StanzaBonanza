package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPoemStore

type mockPoemStore struct{ mock.Mock }

func (m *mockPoemStore) Create(ctx context.Context, poem *domain.Poem) error {
	return m.Called(ctx, poem).Error(0)
}
func (m *mockPoemStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Poem, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Poem)
	return p, args.Error(1)
}
func (m *mockPoemStore) List(ctx context.Context, page domain.PaginationParams, format, sort string) ([]domain.Poem, int, error) {
	args := m.Called(ctx, page, format, sort)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemStore) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, userID, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemStore) ListFeed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, userID, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemStore) ListExplore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemStore) ListHallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	args := m.Called(ctx, page)
	poems, _ := args.Get(0).([]domain.Poem)
	return poems, args.Int(1), args.Error(2)
}
func (m *mockPoemStore) Update(ctx context.Context, poem *domain.Poem) error {
	return m.Called(ctx, poem).Error(0)
}
func (m *mockPoemStore) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockPoemStore) IncrementCounter(ctx context.Context, id uuid.UUID, column string, delta int) error {
	return m.Called(ctx, id, column, delta).Error(0)
}

// mockStanzaStore

type mockStanzaStore struct{ mock.Mock }

func (m *mockStanzaStore) Create(ctx context.Context, stanza *domain.Stanza) error {
	return m.Called(ctx, stanza).Error(0)
}
func (m *mockStanzaStore) ListByPoem(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error) {
	args := m.Called(ctx, poemID)
	stanzas, _ := args.Get(0).([]domain.Stanza)
	return stanzas, args.Error(1)
}
func (m *mockStanzaStore) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.StanzaStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *mockStanzaStore) GetNextPosition(ctx context.Context, poemID uuid.UUID) (int, error) {
	args := m.Called(ctx, poemID)
	return args.Int(0), args.Error(1)
}

// mockPoemNotifStore

type mockPoemNotifStore struct{ mock.Mock }

func (m *mockPoemNotifStore) Create(ctx context.Context, notif *domain.Notification) error {
	return m.Called(ctx, notif).Error(0)
}

func newPoemSvc(poems *mockPoemStore, stanzas *mockStanzaStore, notifs *mockPoemNotifStore) *PoemService {
	return &PoemService{poems: poems, stanzas: stanzas, notifs: notifs}
}

// Create tests

func TestPoemService_Create_Success(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	notifs := &mockPoemNotifStore{}
	svc := newPoemSvc(poems, stanzas, notifs)

	userID := uuid.New()
	poems.On("Create", mock.Anything, mock.AnythingOfType("*domain.Poem")).Return(nil)

	poem, err := svc.Create(context.Background(), userID, "Test Poem", "a poem", domain.FormatFreeVerse, domain.ApprovalOpen, nil)
	require.NoError(t, err)
	assert.Equal(t, "Test Poem", poem.Title)
	assert.Equal(t, userID, poem.AuthorID)
	poems.AssertExpectations(t)
}

func TestPoemService_Create_StoreError(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	poems.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	_, err := svc.Create(context.Background(), uuid.New(), "Title", "", domain.FormatHaiku, domain.ApprovalOpen, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "creating poem")
}

// Get tests

func TestPoemService_Get_Success(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	svc := newPoemSvc(poems, stanzas, &mockPoemNotifStore{})

	poemID := uuid.New()
	expected := &domain.Poem{ID: poemID, Title: "My Poem"}
	stanzaList := []domain.Stanza{{ID: uuid.New(), PoemID: poemID}}

	poems.On("GetByID", mock.Anything, poemID).Return(expected, nil)
	stanzas.On("ListByPoem", mock.Anything, poemID).Return(stanzaList, nil)

	poem, err := svc.Get(context.Background(), poemID)
	require.NoError(t, err)
	assert.Equal(t, "My Poem", poem.Title)
	assert.Len(t, poem.Stanzas, 1)
}

func TestPoemService_Get_NotFound(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	poems.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("no rows"))

	_, err := svc.Get(context.Background(), uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "poem not found")
}

// SubmitStanza tests

func TestPoemService_SubmitStanza_OpenMode_AutoApproved(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	notifs := &mockPoemNotifStore{}
	svc := newPoemSvc(poems, stanzas, notifs)

	authorID := uuid.New()
	submitterID := uuid.New()
	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: authorID, ApprovalMode: domain.ApprovalOpen}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	stanzas.On("GetNextPosition", mock.Anything, poemID).Return(1, nil)
	stanzas.On("Create", mock.Anything, mock.AnythingOfType("*domain.Stanza")).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "stanza_count", 1).Return(nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	stanza, err := svc.SubmitStanza(context.Background(), submitterID, poemID, "Hello world", "")
	require.NoError(t, err)
	assert.Equal(t, domain.StanzaApproved, stanza.Status)
	poems.AssertExpectations(t)
	stanzas.AssertExpectations(t)
}

func TestPoemService_SubmitStanza_ApprovalRequired_Pending(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	notifs := &mockPoemNotifStore{}
	svc := newPoemSvc(poems, stanzas, notifs)

	authorID := uuid.New()
	submitterID := uuid.New()
	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: authorID, ApprovalMode: domain.ApprovalRequired}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	stanzas.On("GetNextPosition", mock.Anything, poemID).Return(1, nil)
	stanzas.On("Create", mock.Anything, mock.AnythingOfType("*domain.Stanza")).Return(nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	stanza, err := svc.SubmitStanza(context.Background(), submitterID, poemID, "text", "")
	require.NoError(t, err)
	assert.Equal(t, domain.StanzaPending, stanza.Status)
}

func TestPoemService_SubmitStanza_Closed_NonAuthor_Error(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	authorID := uuid.New()
	submitterID := uuid.New()
	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: authorID, ApprovalMode: domain.ApprovalClosed}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)

	_, err := svc.SubmitStanza(context.Background(), submitterID, poemID, "text", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestPoemService_SubmitStanza_MaxStanzasReached(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	max := 2
	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: uuid.New(), ApprovalMode: domain.ApprovalOpen, MaxStanzas: &max, StanzaCount: 2}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)

	_, err := svc.SubmitStanza(context.Background(), uuid.New(), poemID, "text", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maximum number of stanzas")
}

func TestPoemService_SubmitStanza_SelfSubmission_AutoApproved(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	notifs := &mockPoemNotifStore{}
	svc := newPoemSvc(poems, stanzas, notifs)

	authorID := uuid.New()
	poemID := uuid.New()
	// Author submits to their own poem which has ApprovalRequired - should still be approved.
	poem := &domain.Poem{ID: poemID, AuthorID: authorID, ApprovalMode: domain.ApprovalRequired}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	stanzas.On("GetNextPosition", mock.Anything, poemID).Return(1, nil)
	stanzas.On("Create", mock.Anything, mock.AnythingOfType("*domain.Stanza")).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "stanza_count", 1).Return(nil)

	stanza, err := svc.SubmitStanza(context.Background(), authorID, poemID, "author text", "")
	require.NoError(t, err)
	assert.Equal(t, domain.StanzaApproved, stanza.Status)
}

// ReviewStanza tests

func TestPoemService_ReviewStanza_Approve_SendsNotification(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	notifs := &mockPoemNotifStore{}
	svc := newPoemSvc(poems, stanzas, notifs)

	authorID := uuid.New()
	submitterID := uuid.New()
	poemID := uuid.New()
	stanzaID := uuid.New()

	poem := &domain.Poem{ID: poemID, AuthorID: authorID}
	stanzaList := []domain.Stanza{{ID: stanzaID, AuthorID: submitterID}}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	stanzas.On("UpdateStatus", mock.Anything, stanzaID, domain.StanzaApproved).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "stanza_count", 1).Return(nil)
	stanzas.On("ListByPoem", mock.Anything, poemID).Return(stanzaList, nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	err := svc.ReviewStanza(context.Background(), authorID, poemID, stanzaID, true)
	require.NoError(t, err)
	notifs.AssertNumberOfCalls(t, "Create", 1)
}

func TestPoemService_ReviewStanza_Reject_SendsNotification(t *testing.T) {
	poems := &mockPoemStore{}
	stanzas := &mockStanzaStore{}
	notifs := &mockPoemNotifStore{}
	svc := newPoemSvc(poems, stanzas, notifs)

	authorID := uuid.New()
	submitterID := uuid.New()
	poemID := uuid.New()
	stanzaID := uuid.New()

	poem := &domain.Poem{ID: poemID, AuthorID: authorID}
	stanzaList := []domain.Stanza{{ID: stanzaID, AuthorID: submitterID}}

	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	stanzas.On("UpdateStatus", mock.Anything, stanzaID, domain.StanzaRejected).Return(nil)
	stanzas.On("ListByPoem", mock.Anything, poemID).Return(stanzaList, nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	err := svc.ReviewStanza(context.Background(), authorID, poemID, stanzaID, false)
	require.NoError(t, err)
	notifs.AssertNumberOfCalls(t, "Create", 1)
}

func TestPoemService_ReviewStanza_NotAuthor_Error(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	authorID := uuid.New()
	otherID := uuid.New()
	poemID := uuid.New()

	poem := &domain.Poem{ID: poemID, AuthorID: authorID}
	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)

	err := svc.ReviewStanza(context.Background(), otherID, poemID, uuid.New(), true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the poem author")
}

// Feed and ListByUser tests

func TestPoemService_Feed_ReturnsPaginatedPoems(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	userID := uuid.New()
	page := domain.PaginationParams{Page: 1, PageSize: 10}
	expected := []domain.Poem{{ID: uuid.New(), Title: "Feed Poem"}}

	poems.On("ListFeed", mock.Anything, userID, page).Return(expected, 1, nil)

	result, total, err := svc.Feed(context.Background(), userID, page)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, result, 1)
}

func TestPoemService_ListByUser_ReturnsPoems(t *testing.T) {
	poems := &mockPoemStore{}
	svc := newPoemSvc(poems, &mockStanzaStore{}, &mockPoemNotifStore{})

	userID := uuid.New()
	page := domain.PaginationParams{Page: 1, PageSize: 20}
	expected := []domain.Poem{{ID: uuid.New(), AuthorID: userID}}

	poems.On("ListByUser", mock.Anything, userID, page).Return(expected, 1, nil)

	result, total, err := svc.ListByUser(context.Background(), userID, page)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, userID, result[0].AuthorID)
}
