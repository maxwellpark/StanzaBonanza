package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/service"
)

type tutorialService interface {
	List(ctx context.Context) ([]domain.Tutorial, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tutorial, error)
}

type TutorialHandler struct {
	svc tutorialService
}

func NewTutorialHandler(svc *service.TutorialService) *TutorialHandler {
	return &TutorialHandler{svc: svc}
}

func (h *TutorialHandler) List(w http.ResponseWriter, r *http.Request) {
	tutorials, err := h.svc.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list tutorials")
		return
	}
	respondJSON(w, http.StatusOK, tutorials)
}

func (h *TutorialHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		respondError(w, http.StatusBadRequest, "slug is required")
		return
	}

	tutorial, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		respondError(w, http.StatusNotFound, "tutorial not found")
		return
	}

	respondJSON(w, http.StatusOK, tutorial)
}
