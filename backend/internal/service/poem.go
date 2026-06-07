package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

type poemStore interface {
	Create(ctx context.Context, poem *domain.Poem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Poem, error)
	List(ctx context.Context, page domain.PaginationParams, format string, sort string) ([]domain.Poem, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error)
	ListFeed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error)
	ListExplore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error)
	ListHallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error)
	Update(ctx context.Context, poem *domain.Poem) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementCounter(ctx context.Context, id uuid.UUID, column string, delta int) error
}

type stanzaStore interface {
	Create(ctx context.Context, stanza *domain.Stanza) error
	ListByPoem(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.StanzaStatus) error
	GetNextPosition(ctx context.Context, poemID uuid.UUID) (int, error)
}

type poemNotifStore interface {
	Create(ctx context.Context, notif *domain.Notification) error
}

type PoemService struct {
	poems   poemStore
	stanzas stanzaStore
	notifs  poemNotifStore
}

func NewPoemService(poems *repository.PoemRepository, stanzas *repository.StanzaRepository, notifs *repository.NotificationRepository) *PoemService {
	return &PoemService{
		poems:   poems,
		stanzas: stanzas,
		notifs:  notifs,
	}
}

func (s *PoemService) Create(ctx context.Context, userID uuid.UUID, title, description string, format domain.PoemFormat, approvalMode domain.ApprovalMode, maxStanzas *int) (*domain.Poem, error) {
	var poem = &domain.Poem{
		AuthorID:     userID,
		Title:        title,
		Description:  description,
		Format:       format,
		ApprovalMode: approvalMode,
		MaxStanzas:   maxStanzas,
	}

	if err := s.poems.Create(ctx, poem); err != nil {
		return nil, fmt.Errorf("creating poem: %w", err)
	}

	return poem, nil
}

func (s *PoemService) Get(ctx context.Context, id uuid.UUID) (*domain.Poem, error) {
	poem, err := s.poems.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("poem not found: %w", err)
	}

	stanzas, err := s.stanzas.ListByPoem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("loading stanzas: %w", err)
	}
	poem.Stanzas = stanzas

	return poem, nil
}

func (s *PoemService) List(ctx context.Context, page domain.PaginationParams, format, sort string) ([]domain.Poem, int, error) {
	return s.poems.List(ctx, page, format, sort)
}

func (s *PoemService) ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	return s.poems.ListByUser(ctx, userID, page)
}

func (s *PoemService) Feed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error) {
	return s.poems.ListFeed(ctx, userID, page)
}

func (s *PoemService) Explore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	return s.poems.ListExplore(ctx, page)
}

func (s *PoemService) HallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error) {
	return s.poems.ListHallOfFame(ctx, page)
}

func (s *PoemService) Update(ctx context.Context, userID, poemID uuid.UUID, title, description string) error {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return fmt.Errorf("poem not found: %w", err)
	}
	if poem.AuthorID != userID {
		return fmt.Errorf("not the poem author")
	}

	poem.Title = title
	poem.Description = description
	return s.poems.Update(ctx, poem)
}

func (s *PoemService) Delete(ctx context.Context, userID, poemID uuid.UUID) error {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return fmt.Errorf("poem not found: %w", err)
	}
	if poem.AuthorID != userID {
		return fmt.Errorf("not the poem author")
	}

	return s.poems.Delete(ctx, poemID)
}

func (s *PoemService) ListStanzas(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error) {
	return s.stanzas.ListByPoem(ctx, poemID)
}

func (s *PoemService) SubmitStanza(ctx context.Context, userID, poemID uuid.UUID, text, literaryDevice string) (*domain.Stanza, error) {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return nil, fmt.Errorf("poem not found: %w", err)
	}

	if poem.ApprovalMode == domain.ApprovalClosed && poem.AuthorID != userID {
		return nil, fmt.Errorf("poem is closed for submissions")
	}

	if poem.MaxStanzas != nil && poem.StanzaCount >= *poem.MaxStanzas {
		return nil, fmt.Errorf("poem has reached the maximum number of stanzas")
	}

	nextPos, err := s.stanzas.GetNextPosition(ctx, poemID)
	if err != nil {
		return nil, fmt.Errorf("getting next position: %w", err)
	}

	var status domain.StanzaStatus
	if poem.ApprovalMode == domain.ApprovalRequired && poem.AuthorID != userID {
		status = domain.StanzaPending
	} else {
		status = domain.StanzaApproved
	}

	var stanza = &domain.Stanza{
		PoemID:         poemID,
		AuthorID:       userID,
		Text:           text,
		Position:       nextPos,
		LiteraryDevice: literaryDevice,
		Status:         status,
	}

	if err := s.stanzas.Create(ctx, stanza); err != nil {
		return nil, fmt.Errorf("creating stanza: %w", err)
	}

	if status == domain.StanzaApproved {
		_ = s.poems.IncrementCounter(ctx, poemID, "stanza_count", 1)
	}

	if poem.AuthorID != userID {
		_ = s.notifs.Create(ctx, &domain.Notification{
			RecipientID: poem.AuthorID,
			ActorID:     &userID,
			Type:        domain.NotifStanzaSubmit,
			PoemID:      &poemID,
		})
	}

	return stanza, nil
}

func (s *PoemService) ReviewStanza(ctx context.Context, userID, poemID, stanzaID uuid.UUID, approved bool) error {
	poem, err := s.poems.GetByID(ctx, poemID)
	if err != nil {
		return fmt.Errorf("poem not found: %w", err)
	}
	if poem.AuthorID != userID {
		return fmt.Errorf("not the poem author")
	}

	var status domain.StanzaStatus
	var notifType domain.NotificationType
	if approved {
		status = domain.StanzaApproved
		notifType = domain.NotifStanzaApproved
		_ = s.poems.IncrementCounter(ctx, poemID, "stanza_count", 1)
	} else {
		status = domain.StanzaRejected
		notifType = domain.NotifStanzaRejected
	}

	if err := s.stanzas.UpdateStatus(ctx, stanzaID, status); err != nil {
		return fmt.Errorf("updating stanza status: %w", err)
	}

	stanzas, err := s.stanzas.ListByPoem(ctx, poemID)
	if err == nil {
		for _, st := range stanzas {
			if st.ID == stanzaID && st.AuthorID != userID {
				_ = s.notifs.Create(ctx, &domain.Notification{
					RecipientID: st.AuthorID,
					ActorID:     &userID,
					Type:        notifType,
					PoemID:      &poemID,
				})
				break
			}
		}
	}

	return nil
}
