package handler

import (
	"context"
	"net/http"
	"strings"

	"asenare/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "userID"

type UserHandler struct {
	users repository.UserRepository
}

func NewUserHandler(users repository.UserRepository) *UserHandler {
	return &UserHandler{users: users}
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
