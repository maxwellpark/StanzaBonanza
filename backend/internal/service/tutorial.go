package service

import (
	"context"
	"fmt"

	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/repository"
)

type tutorialStore interface {
	List(ctx context.Context) ([]domain.Tutorial, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tutorial, error)
}

type TutorialService struct {
	repo tutorialStore
}

func NewTutorialService(repo *repository.TutorialRepository) *TutorialService {
	return &TutorialService{repo: repo}
}

func (s *TutorialService) List(ctx context.Context) ([]domain.Tutorial, error) {
	tutorials, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing tutorials: %w", err)
	}
	if tutorials == nil {
		tutorials = []domain.Tutorial{}
	}
	return tutorials, nil
}

func (s *TutorialService) GetBySlug(ctx context.Context, slug string) (*domain.Tutorial, error) {
	tutorial, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("tutorial not found: %w", err)
	}
	return tutorial, nil
}
