package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/maxwellpark/stanzabonanza/backend/internal/service"
)

type SocialHandler struct {
	svc *service.SocialService
}

func NewSocialHandler(svc *service.SocialService) *SocialHandler {
	return &SocialHandler{svc: svc}
}

func (h *SocialHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
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

	liked, err := h.svc.ToggleLike(r.Context(), userID, poemID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to toggle like")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

func (h *SocialHandler) AddComment(w http.ResponseWriter, r *http.Request) {
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
		ParentID *uuid.UUID `json:"parentId"`
		Text     string     `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Text == "" {
		respondError(w, http.StatusBadRequest, "text is required")
		return
	}
	if len(body.Text) > 2000 {
		respondError(w, http.StatusBadRequest, "comment must be 2000 characters or fewer")
		return
	}

	comment, err := h.svc.AddComment(r.Context(), userID, poemID, body.ParentID, body.Text)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add comment")
		return
	}

	respondJSON(w, http.StatusCreated, comment)
}

func (h *SocialHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid comment ID")
		return
	}

	if err := h.svc.DeleteComment(r.Context(), userID, commentID); err != nil {
		if err.Error() == "not the comment author" {
			respondError(w, http.StatusForbidden, "you are not the author of this comment")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete comment")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "comment deleted"})
}

func (h *SocialHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	poemID, err := uuid.Parse(chi.URLParam(r, "poemID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid poem ID")
		return
	}

	var page = parsePagination(r)
	comments, total, err := h.svc.ListComments(r.Context(), poemID, page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list comments")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Comment]{
		Items:      comments,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *SocialHandler) ToggleFollow(w http.ResponseWriter, r *http.Request) {
	followerID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	followedID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	following, err := h.svc.ToggleFollow(r.Context(), followerID, followedID)
	if err != nil {
		if err.Error() == "cannot follow yourself" {
			respondError(w, http.StatusBadRequest, "cannot follow yourself")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to toggle follow")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"following": following})
}

func (h *SocialHandler) ListFollowers(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var page = parsePagination(r)
	users, total, err := h.svc.ListFollowers(r.Context(), userID, page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list followers")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.User]{
		Items:      users,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *SocialHandler) ListFollowing(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var page = parsePagination(r)
	users, total, err := h.svc.ListFollowing(r.Context(), userID, page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list following")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.User]{
		Items:      users,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *SocialHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var page = parsePagination(r)
	notifs, total, err := h.svc.ListNotifications(r.Context(), userID, page)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}

	respondJSON(w, http.StatusOK, domain.PaginatedResult[domain.Notification]{
		Items:      notifs,
		TotalCount: total,
		Page:       page.Page,
		PageSize:   page.PageSize,
	})
}

func (h *SocialHandler) MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		IDs []uuid.UUID `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(body.IDs) == 0 {
		respondError(w, http.StatusBadRequest, "ids is required")
		return
	}

	if err := h.svc.MarkNotificationsRead(r.Context(), userID, body.IDs); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark notifications read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "notifications marked as read"})
}
