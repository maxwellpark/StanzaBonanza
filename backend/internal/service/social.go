package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

type SocialService struct {
	likes    *repository.LikeRepository
	comments *repository.CommentRepository
	follows  *repository.FollowRepository
	notifs   *repository.NotificationRepository
	poems    *repository.PoemRepository
}

func NewSocialService(
	likes *repository.LikeRepository,
	comments *repository.CommentRepository,
	follows *repository.FollowRepository,
	notifs *repository.NotificationRepository,
	poems *repository.PoemRepository,
) *SocialService {
	return &SocialService{
		likes:    likes,
		comments: comments,
		follows:  follows,
		notifs:   notifs,
		poems:    poems,
	}
}

func (s *SocialService) ToggleLike(ctx context.Context, userID, poemID uuid.UUID) (bool, error) {
	exists, err := s.likes.Exists(ctx, userID, poemID)
	if err != nil {
		return false, fmt.Errorf("checking like: %w", err)
	}

	if exists {
		if err := s.likes.Delete(ctx, userID, poemID); err != nil {
			return false, fmt.Errorf("removing like: %w", err)
		}
		_ = s.poems.IncrementCounter(ctx, poemID, "like_count", -1)
		return false, nil
	}

	var like = &domain.Like{
		UserID: userID,
		PoemID: poemID,
	}
	if err := s.likes.Create(ctx, like); err != nil {
		return false, fmt.Errorf("adding like: %w", err)
	}
	_ = s.poems.IncrementCounter(ctx, poemID, "like_count", 1)

	poem, err := s.poems.GetByID(ctx, poemID)
	if err == nil && poem.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: poem.AuthorID,
			ActorID:     &userID,
			Type:        domain.NotifLike,
			PoemID:      &poemID,
		})
	}

	return true, nil
}

func (s *SocialService) AddComment(ctx context.Context, userID, poemID uuid.UUID, parentID *uuid.UUID, text string) (*domain.Comment, error) {
	var comment = &domain.Comment{
		PoemID:   poemID,
		AuthorID: userID,
		ParentID: parentID,
		Text:     text,
	}

	if err := s.comments.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}

	_ = s.poems.IncrementCounter(ctx, poemID, "comment_count", 1)

	poem, err := s.poems.GetByID(ctx, poemID)
	if err == nil && poem.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: poem.AuthorID,
			ActorID:     &userID,
			Type:        domain.NotifComment,
			PoemID:      &poemID,
		})
	}

	return comment, nil
}

func (s *SocialService) DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error {
	comment, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}
	if comment.AuthorID != userID {
		return fmt.Errorf("not the comment author")
	}

	if err := s.comments.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	_ = s.poems.IncrementCounter(ctx, comment.PoemID, "comment_count", -1)
	return nil
}

func (s *SocialService) ListComments(ctx context.Context, poemID uuid.UUID, page domain.PaginationParams) ([]domain.Comment, int, error) {
	return s.comments.ListByPoem(ctx, poemID, page)
}

func (s *SocialService) ToggleFollow(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	if followerID == followedID {
		return false, fmt.Errorf("cannot follow yourself")
	}

	exists, err := s.follows.Exists(ctx, followerID, followedID)
	if err != nil {
		return false, fmt.Errorf("checking follow: %w", err)
	}

	if exists {
		if err := s.follows.Delete(ctx, followerID, followedID); err != nil {
			return false, fmt.Errorf("unfollowing: %w", err)
		}
		return false, nil
	}

	var follow = &domain.Follow{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	if err := s.follows.Create(ctx, follow); err != nil {
		return false, fmt.Errorf("following: %w", err)
	}

	_ = s.notifs.Create(ctx, &domain.Notification{
		RecipientID: followedID,
		ActorID:     &followerID,
		Type:        domain.NotifFollow,
	})

	return true, nil
}

func (s *SocialService) ListFollowers(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	return s.follows.ListFollowers(ctx, userID, page)
}

func (s *SocialService) ListFollowing(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.User, int, error) {
	return s.follows.ListFollowing(ctx, userID, page)
}

func (s *SocialService) ListNotifications(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Notification, int, error) {
	return s.notifs.ListByUser(ctx, userID, page)
}

func (s *SocialService) MarkNotificationsRead(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	return s.notifs.MarkRead(ctx, userID, ids)
}
