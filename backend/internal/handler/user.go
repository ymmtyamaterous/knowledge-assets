package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"asenare/backend/internal/repository"
	"asenare/backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "userID"

type UserHandler struct {
	users  repository.UserRepository
	authUC *usecase.AuthUseCase
}

func NewUserHandler(users repository.UserRepository, authUC *usecase.AuthUseCase) *UserHandler {
	return &UserHandler{users: users, authUC: authUC}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, ok, err := h.users.FindByID(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !ok {
		WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	WriteJSON(w, http.StatusOK, user)
}

type updateMeRequest struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, ok, err := h.users.FindByID(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !ok {
		WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	var req updateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username != "" {
		user.Username = strings.TrimSpace(req.Username)
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	updated, err := h.users.Update(user)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, updated)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		WriteError(w, http.StatusBadRequest, "currentPassword and newPassword are required")
		return
	}

	if err := h.authUC.ChangePassword(userID, req.CurrentPassword, req.NewPassword); err != nil {
		if err.Error() == usecase.ErrInvalidCredentials.Error() {
			WriteError(w, http.StatusBadRequest, "current password is incorrect or new password is too short")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func JWTAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(auth, "Bearer ") {
				WriteError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			rawToken := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
			tok, err := jwt.Parse(rawToken, func(token *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !tok.Valid {
				WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := tok.Claims.(jwt.MapClaims)
			if !ok {
				WriteError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			sub, _ := claims["sub"].(string)
			if sub == "" {
				WriteError(w, http.StatusUnauthorized, "invalid subject")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
