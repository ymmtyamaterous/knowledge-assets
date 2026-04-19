package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"asenare/backend/internal/repository"
	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

func TestAuthAndMe(t *testing.T) {
	userRepo := repository.NewInMemoryUserRepository()
	authUC := usecase.NewAuthUseCase(userRepo, "test-secret")

	authHandler := NewAuthHandler(authUC)
	userHandler := NewUserHandler(userRepo, authUC)

	r := chi.NewRouter()
	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Group(func(private chi.Router) {
		private.Use(JWTAuthMiddleware("test-secret"))
		private.Get("/api/v1/users/me", userHandler.Me)
	})

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"username": "tester",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var registerResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&registerResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	if registerResp.Token == "" {
		t.Fatal("token is empty")
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+registerResp.Token)
	meRec := httptest.NewRecorder()
	r.ServeHTTP(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", meRec.Code)
	}
}
