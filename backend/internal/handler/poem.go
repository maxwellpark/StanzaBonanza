package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/maxwellpark/stanzabonanza/backend/internal/service"
)

type poemService interface {
	Create(ctx context.Context, userID uuid.UUID, title, description string, format domain.PoemFormat, approvalMode domain.ApprovalMode, maxStanzas *int) (*domain.Poem, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Poem, error)
	List(ctx context.Context, page domain.PaginationParams, format, sort string) ([]domain.Poem, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error)
	Update(ctx context.Context, userID, poemID uuid.UUID, title, description string) error
	Delete(ctx context.Context, userID, poemID uuid.UUID) error
	ListStanzas(ctx context.Context, poemID uuid.UUID) ([]domain.Stanza, error)
	SubmitStanza(ctx context.Context, userID, poemID uuid.UUID, text, literaryDevice string) (*domain.Stanza, error)
	ReviewStanza(ctx context.Context, userID, poemID, stanzaID uuid.UUID, approved bool) error
	Feed(ctx context.Context, userID uuid.UUID, page domain.PaginationParams) ([]domain.Poem, int, error)
	Explore(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error)
	HallOfFame(ctx context.Context, page domain.PaginationParams) ([]domain.Poem, int, error)
}

type PoemHandler struct {
	svc poemService
}

func NewPoemHandler(svc *service.PoemService) *PoemHandler {
	return &PoemHandler{svc: svc}
}

func (h *PoemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		Title        string             `json:"title"`
		Description  string             `json:"description"`
		Format       domain.PoemFormat  `json:"format"`
		ApprovalMode domain.ApprovalMode `json:"approvalMode"`
		MaxStanzas   *int               `json:"maxStanzas"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}
	if len(body.Title) > 200 {
		respondError(w, http.StatusBadRequest, "title must be 200 characters or fewer")
		return
	}
	if body.Format == "" {
		respondError(w, http.StatusBadRequest, "format is required")
		return
	}
	if body.ApprovalMode == "" {
		body.ApprovalMode = domain.ApprovalOpen
	}

	poem, err := h.svc.Create(r.Context(), userID, body.Title, body.Description, body.Format, body.ApprovalMode, body.MaxStanzas)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create poem")
		return
	}

	respondJSON(w, http.StatusCreated, poem)
}

func (h *PoemHandler) Get(w http.ResponseWriter, r *http.Request) {
	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	poem, err := h.svc.Get(r.Context(), poemID)
	if err != nil {
		respondError(w, http.StatusNotFound, "poem not found")
		return
	}

	respondJSON(w, http.StatusOK, poem)
}

func (h *PoemHandler) List(w http.ResponseWriter, r *http.Request) {
	var page = parsePagination(r)
	var format = r.URL.Query().Get("format")
	var sort = r.URL.Query().Get("sort")

	poems, total, err := h.svc.List(r.Context(), page, format, sort)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list poems")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Poem]{
		Items:      poems,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *PoemHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var page = parsePagination(r)
	poems, total, err := h.svc.ListByUser(r.Context(), userID, page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list poems")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Poem]{
		Items:      poems,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *PoemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}
	if len(body.Title) > 200 {
		respondError(w, http.StatusBadRequest, "title must be 200 characters or fewer")
		return
	}

	if err := h.svc.Update(r.Context(), userID, poemID, body.Title, body.Description); err != nil {
		if err.Error() == "not the poem author" {
			respondError(w, http.StatusForbidden, "you are not the author of this poem")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update poem")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "poem updated"})
}

func (h *PoemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	if err := h.svc.Delete(r.Context(), userID, poemID); err != nil {
		if err.Error() == "not the poem author" {
			respondError(w, http.StatusForbidden, "you are not the author of this poem")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete poem")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "poem deleted"})
}

func (h *PoemHandler) ListStanzas(w http.ResponseWriter, r *http.Request) {
	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	stanzas, err := h.svc.ListStanzas(r.Context(), poemID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list stanzas")
		return
	}

	respondJSON(w, http.StatusOK, stanzas)
}

func (h *PoemHandler) SubmitStanza(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	var body struct {
		Text           string `json:"text"`
		LiteraryDevice string `json:"literaryDevice"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Text == "" {
		respondError(w, http.StatusBadRequest, "text is required")
		return
	}
	if len(body.Text) > 10000 {
		respondError(w, http.StatusBadRequest, "text must be 10000 characters or fewer")
		return
	}

	stanza, err := h.svc.SubmitStanza(r.Context(), userID, poemID, body.Text, body.LiteraryDevice)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, stanza)
}

func (h *PoemHandler) ReviewStanza(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	stanzaID, err := uuid.Parse(chi.URLParam(r, "stanzaID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stanza ID")
		return
	}

	var body struct {
		Approved bool `json:"approved"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.ReviewStanza(r.Context(), userID, poemID, stanzaID, body.Approved); err != nil {
		if err.Error() == "not the poem author" {
			respondError(w, http.StatusForbidden, "you are not the author of this poem")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to review stanza")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "stanza reviewed"})
}

func (h *PoemHandler) Feed(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var page = parsePagination(r)
	poems, total, err := h.svc.Feed(r.Context(), userID, page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load feed")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Poem]{
		Items:      poems,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *PoemHandler) Explore(w http.ResponseWriter, r *http.Request) {
	var page = parsePagination(r)
	poems, total, err := h.svc.Explore(r.Context(), page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load explore")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Poem]{
		Items:      poems,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *PoemHandler) HallOfFame(w http.ResponseWriter, r *http.Request) {
	var page = parsePagination(r)
	poems, total, err := h.svc.HallOfFame(r.Context(), page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load hall of fame")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Poem]{
		Items:      poems,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}
