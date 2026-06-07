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

type mockTutorialStore struct{ mock.Mock }

func (m *mockTutorialStore) List(ctx context.Context) ([]domain.Tutorial, error) {
	args := m.Called(ctx)
	t, _ := args.Get(0).([]domain.Tutorial)
	return t, args.Error(1)
}
func (m *mockTutorialStore) GetBySlug(ctx context.Context, slug string) (*domain.Tutorial, error) {
	args := m.Called(ctx, slug)
	t, _ := args.Get(0).(*domain.Tutorial)
	return t, args.Error(1)
}

func TestTutorialService_List_ReturnsTutorials(t *testing.T) {
	store := &mockTutorialStore{}
	svc := &TutorialService{repo: store}

	expected := []domain.Tutorial{
		{ID: uuid.New(), Title: "Haiku Basics", Slug: "haiku-basics"},
	}
	store.On("List", mock.Anything).Return(expected, nil)

	result, err := svc.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Haiku Basics", result[0].Title)
}

func TestTutorialService_List_EmptySlice_NotNil(t *testing.T) {
	store := &mockTutorialStore{}
	svc := &TutorialService{repo: store}

	// Store returns nil - service should convert to empty slice.
	store.On("List", mock.Anything).Return(nil, nil)

	result, err := svc.List(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestTutorialService_List_Error(t *testing.T) {
	store := &mockTutorialStore{}
	svc := &TutorialService{repo: store}

	store.On("List", mock.Anything).Return(nil, errors.New("db error"))

	_, err := svc.List(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "listing tutorials")
}

func TestTutorialService_GetBySlug_Success(t *testing.T) {
	store := &mockTutorialStore{}
	svc := &TutorialService{repo: store}

	expected := &domain.Tutorial{ID: uuid.New(), Slug: "haiku-basics", Title: "Haiku Basics"}
	store.On("GetBySlug", mock.Anything, "haiku-basics").Return(expected, nil)

	result, err := svc.GetBySlug(context.Background(), "haiku-basics")
	require.NoError(t, err)
	assert.Equal(t, "haiku-basics", result.Slug)
}

func TestTutorialService_GetBySlug_NotFound(t *testing.T) {
	store := &mockTutorialStore{}
	svc := &TutorialService{repo: store}

	store.On("GetBySlug", mock.Anything, "missing").Return(nil, errors.New("not found"))

	_, err := svc.GetBySlug(context.Background(), "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tutorial not found")
}
