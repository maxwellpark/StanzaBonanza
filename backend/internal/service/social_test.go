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

// mockLikeStore

type mockLikeStore struct{ mock.Mock }

func (m *mockLikeStore) Create(ctx context.Context, like *domain.Like) error {
	return m.Called(ctx, like).Error(0)
}
func (m *mockLikeStore) Delete(ctx context.Context, userID, poemID uuid.UUID) error {
	return m.Called(ctx, userID, poemID).Error(0)
}
func (m *mockLikeStore) Exists(ctx context.Context, userID, poemID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, poemID)
	return args.Bool(0), args.Error(1)
}

// mockCommentStore

type mockCommentStore struct{ mock.Mock }

func (m *mockCommentStore) Create(ctx context.Context, comment *domain.Comment) error {
	return m.Called(ctx, comment).Error(0)
}
func (m *mockCommentStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*domain.Comment)
	return c, args.Error(1)
}
func (m *mockCommentStore) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockCommentStore) ListByPoem(ctx context.Context, poemID uuid.UUID, page domain.PaginationParams) ([]domain.Comment, int, error) {
	args := m.Called(ctx, poemID, page)
	c, _ := args.Get(0).([]domain.Comment)
	return c, args.Int(1), args.Error(2)
}

// mockFollowStore

type mockFollowStore struct{ mock.Mock }

func (m *mockFollowStore) Create(ctx context.Context, follow *domain.Follow) error {
	return m.Called(ctx, follow).Error(0)
}
func (m *mockFollowStore) Delete(ctx context.Context, followerID, followedID uuid.UUID) error {
	return m.Called(ctx, followerID, followedID).Error(0)
}
func (m *mockFollowStore) Exists(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0), args.Error(1)
}
func (m *mockFollowStore) ListFollowers(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	args := m.Called(ctx, userID, page)
	u, _ := args.Get(0).([]domain.User)
	return u, args.Int(1), args.Error(2)
}
func (m *mockFollowStore) ListFollowing(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	args := m.Called(ctx, userID, page)
	u, _ := args.Get(0).([]domain.User)
	return u, args.Int(1), args.Error(2)
}

// mockSocialNotifStore

type mockSocialNotifStore struct{ mock.Mock }

func (m *mockSocialNotifStore) Create(ctx context.Context, notif *domain.Notification) error {
	return m.Called(ctx, notif).Error(0)
}
func (m *mockSocialNotifStore) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Notification, int, error) {
	args := m.Called(ctx, userID, page)
	n, _ := args.Get(0).([]domain.Notification)
	return n, args.Int(1), args.Error(2)
}
func (m *mockSocialNotifStore) MarkRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return m.Called(ctx, userID, ids).Error(0)
}

// mockSocialPoemStore

type mockSocialPoemStore struct{ mock.Mock }

func (m *mockSocialPoemStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Poem, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Poem)
	return p, args.Error(1)
}
func (m *mockSocialPoemStore) IncrementCounter(ctx context.Context, id uuid.UUID, column string, delta int) error {
	return m.Called(ctx, id, column, delta).Error(0)
}

func newSocialSvc(likes *mockLikeStore, comments *mockCommentStore, follows *mockFollowStore, notifs *mockSocialNotifStore, poems *mockSocialPoemStore) *SocialService {
	return &SocialService{
		likes:    likes,
		comments: comments,
		follows:  follows,
		notifs:   notifs,
		poems:    poems,
	}
}

// ToggleLike tests

func TestSocialService_ToggleLike_Like(t *testing.T) {
	likes := &mockLikeStore{}
	poems := &mockSocialPoemStore{}
	notifs := &mockSocialNotifStore{}
	svc := newSocialSvc(likes, &mockCommentStore{}, &mockFollowStore{}, notifs, poems)

	userID := uuid.New()
	authorID := uuid.New()
	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: authorID}

	likes.On("Exists", mock.Anything, userID, poemID).Return(false, nil)
	likes.On("Create", mock.Anything, mock.AnythingOfType("*domain.Like")).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "like_count", 1).Return(nil)
	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	liked, err := svc.ToggleLike(context.Background(), userID, poemID)
	require.NoError(t, err)
	assert.True(t, liked)
	notifs.AssertNumberOfCalls(t, "Create", 1)
}

func TestSocialService_ToggleLike_Unlike(t *testing.T) {
	likes := &mockLikeStore{}
	poems := &mockSocialPoemStore{}
	svc := newSocialSvc(likes, &mockCommentStore{}, &mockFollowStore{}, &mockSocialNotifStore{}, poems)

	userID := uuid.New()
	poemID := uuid.New()

	likes.On("Exists", mock.Anything, userID, poemID).Return(true, nil)
	likes.On("Delete", mock.Anything, userID, poemID).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "like_count", -1).Return(nil)

	liked, err := svc.ToggleLike(context.Background(), userID, poemID)
	require.NoError(t, err)
	assert.False(t, liked)
}

// AddComment tests

func TestSocialService_AddComment_Success(t *testing.T) {
	comments := &mockCommentStore{}
	poems := &mockSocialPoemStore{}
	notifs := &mockSocialNotifStore{}
	svc := newSocialSvc(&mockLikeStore{}, comments, &mockFollowStore{}, notifs, poems)

	userID := uuid.New()
	authorID := uuid.New()
	poemID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: authorID}

	comments.On("Create", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "comment_count", 1).Return(nil)
	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	comment, err := svc.AddComment(context.Background(), userID, poemID, nil, "Hello")
	require.NoError(t, err)
	assert.Equal(t, "Hello", comment.Text)
}

func TestSocialService_AddComment_WithParent(t *testing.T) {
	comments := &mockCommentStore{}
	poems := &mockSocialPoemStore{}
	notifs := &mockSocialNotifStore{}
	svc := newSocialSvc(&mockLikeStore{}, comments, &mockFollowStore{}, notifs, poems)

	userID := uuid.New()
	authorID := uuid.New()
	poemID := uuid.New()
	parentID := uuid.New()
	poem := &domain.Poem{ID: poemID, AuthorID: authorID}

	comments.On("Create", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "comment_count", 1).Return(nil)
	poems.On("GetByID", mock.Anything, poemID).Return(poem, nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	comment, err := svc.AddComment(context.Background(), userID, poemID, &parentID, "Reply text")
	require.NoError(t, err)
	require.NotNil(t, comment.ParentID)
	assert.Equal(t, parentID, *comment.ParentID)
}

// DeleteComment tests

func TestSocialService_DeleteComment_Success(t *testing.T) {
	comments := &mockCommentStore{}
	poems := &mockSocialPoemStore{}
	svc := newSocialSvc(&mockLikeStore{}, comments, &mockFollowStore{}, &mockSocialNotifStore{}, poems)

	userID := uuid.New()
	commentID := uuid.New()
	poemID := uuid.New()
	comment := &domain.Comment{ID: commentID, AuthorID: userID, PoemID: poemID}

	comments.On("GetByID", mock.Anything, commentID).Return(comment, nil)
	comments.On("Delete", mock.Anything, commentID).Return(nil)
	poems.On("IncrementCounter", mock.Anything, poemID, "comment_count", -1).Return(nil)

	err := svc.DeleteComment(context.Background(), userID, commentID)
	require.NoError(t, err)
}

func TestSocialService_DeleteComment_NotAuthor(t *testing.T) {
	comments := &mockCommentStore{}
	svc := newSocialSvc(&mockLikeStore{}, comments, &mockFollowStore{}, &mockSocialNotifStore{}, &mockSocialPoemStore{})

	userID := uuid.New()
	otherID := uuid.New()
	commentID := uuid.New()
	comment := &domain.Comment{ID: commentID, AuthorID: otherID, PoemID: uuid.New()}

	comments.On("GetByID", mock.Anything, commentID).Return(comment, nil)

	err := svc.DeleteComment(context.Background(), userID, commentID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not the comment author")
}

// ToggleFollow tests

func TestSocialService_ToggleFollow_Follow(t *testing.T) {
	follows := &mockFollowStore{}
	notifs := &mockSocialNotifStore{}
	svc := newSocialSvc(&mockLikeStore{}, &mockCommentStore{}, follows, notifs, &mockSocialPoemStore{})

	followerID := uuid.New()
	followedID := uuid.New()

	follows.On("Exists", mock.Anything, followerID, followedID).Return(false, nil)
	follows.On("Create", mock.Anything, mock.AnythingOfType("*domain.Follow")).Return(nil)
	notifs.On("Create", mock.Anything, mock.AnythingOfType("*domain.Notification")).Return(nil)

	following, err := svc.ToggleFollow(context.Background(), followerID, followedID)
	require.NoError(t, err)
	assert.True(t, following)
	notifs.AssertNumberOfCalls(t, "Create", 1)
}

func TestSocialService_ToggleFollow_Unfollow(t *testing.T) {
	follows := &mockFollowStore{}
	svc := newSocialSvc(&mockLikeStore{}, &mockCommentStore{}, follows, &mockSocialNotifStore{}, &mockSocialPoemStore{})

	followerID := uuid.New()
	followedID := uuid.New()

	follows.On("Exists", mock.Anything, followerID, followedID).Return(true, nil)
	follows.On("Delete", mock.Anything, followerID, followedID).Return(nil)

	following, err := svc.ToggleFollow(context.Background(), followerID, followedID)
	require.NoError(t, err)
	assert.False(t, following)
}

// ListNotifications tests

func TestSocialService_ListNotifications(t *testing.T) {
	notifs := &mockSocialNotifStore{}
	svc := newSocialSvc(&mockLikeStore{}, &mockCommentStore{}, &mockFollowStore{}, notifs, &mockSocialPoemStore{})

	userID := uuid.New()
	page := domain.PaginationParams{Page: 1, PageSize: 20}
	expected := []domain.Notification{{ID: uuid.New(), RecipientID: userID, Type: domain.NotifLike}}

	notifs.On("ListByUser", mock.Anything, userID, page).Return(expected, 1, nil)

	result, total, err := svc.ListNotifications(context.Background(), userID, page)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, result, 1)
}

// MarkNotificationsRead tests

func TestSocialService_MarkNotificationsRead_Success(t *testing.T) {
	notifs := &mockSocialNotifStore{}
	svc := newSocialSvc(&mockLikeStore{}, &mockCommentStore{}, &mockFollowStore{}, notifs, &mockSocialPoemStore{})

	userID := uuid.New()
	ids := []uuid.UUID{uuid.New(), uuid.New()}

	notifs.On("MarkRead", mock.Anything, userID, ids).Return(nil)

	err := svc.MarkNotificationsRead(context.Background(), userID, ids)
	require.NoError(t, err)

	// Error case.
	notifs2 := &mockSocialNotifStore{}
	svc2 := newSocialSvc(&mockLikeStore{}, &mockCommentStore{}, &mockFollowStore{}, notifs2, &mockSocialPoemStore{})
	notifs2.On("MarkRead", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))

	err = svc2.MarkNotificationsRead(context.Background(), userID, ids)
	require.Error(t, err)
}
