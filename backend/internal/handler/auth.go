package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"github.com/maxwellpark/stanzabonanza/backend/internal/domain"
	"github.com/maxwellpark/stanzabonanza/backend/internal/middleware"
	"github.com/maxwellpark/stanzabonanza/backend/internal/service"
)

type authService interface {
	RequestMagicLink(ctx context.Context, email string) (string, error)
	VerifyMagicLink(ctx context.Context, token string) (string, *domain.User, error)
	DeleteSession(ctx context.Context, token string) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, avatarURL string) error
	BeginRegistration(ctx context.Context, userID uuid.UUID) (*protocol.CredentialCreation, string, error)
	FinishRegistration(ctx context.Context, userID uuid.UUID, sessionKey string, r *http.Request) (*domain.WebAuthnCredential, error)
	BeginLogin(ctx context.Context) (*protocol.CredentialAssertion, string, error)
	FinishLogin(ctx context.Context, sessionKey string, r *http.Request) (*domain.User, string, error)
}

type AuthHandler struct {
	svc authService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RequestMagicLink(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Email == "" {
		respondError(w, http.StatusBadRequest, "email is required")
		return
	}

	_, err := h.svc.RequestMagicLink(r.Context(), body.Email)
	if err != nil {
		// Don't leak whether the email exists
		respondJSON(w, http.StatusOK, map[string]string{"message": "if that email is registered, a magic link has been sent"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "if that email is registered, a magic link has been sent"})
}

func (h *AuthHandler) VerifyMagicLink(w http.ResponseWriter, r *http.Request) {
	var token = r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusBadRequest, "token is required")
		return
	}

	sessionToken, user, err := h.svc.VerifyMagicLink(r.Context(), token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid or expired magic link")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60,
	})

	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) BeginRegistration(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	creation, sessionKey, err := h.svc.BeginRegistration(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to begin registration")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "wa_session",
		Value:    sessionKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	respondJSON(w, http.StatusOK, creation)
}

func (h *AuthHandler) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	cookie, err := r.Cookie("wa_session")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing WebAuthn session")
		return
	}

	if _, err := h.svc.FinishRegistration(r.Context(), userID, cookie.Value, r); err != nil {
		respondError(w, http.StatusBadRequest, "registration failed: "+err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "wa_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "passkey registered"})
}

func (h *AuthHandler) BeginLogin(w http.ResponseWriter, r *http.Request) {
	assertion, sessionKey, err := h.svc.BeginLogin(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to begin login")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "wa_session",
		Value:    sessionKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	respondJSON(w, http.StatusOK, assertion)
}

func (h *AuthHandler) FinishLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("wa_session")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing WebAuthn session")
		return
	}

	user, sessionToken, err := h.svc.FinishLogin(r.Context(), cookie.Value, r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "login failed: "+err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "wa_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		_ = h.svc.DeleteSession(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.svc.GetUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	var idStr = chi.URLParam(r, "userID")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.svc.GetUser(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		DisplayName string `json:"displayName"`
		Bio         string `json:"bio"`
		AvatarURL   string `json:"avatarUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(body.DisplayName) > 50 {
		respondError(w, http.StatusBadRequest, "display name must be 50 characters or fewer")
		return
	}

	if err := h.svc.UpdateProfile(r.Context(), userID, body.DisplayName, body.Bio, body.AvatarURL); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "profile updated"})
}
