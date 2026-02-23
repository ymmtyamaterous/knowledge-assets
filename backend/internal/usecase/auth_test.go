package usecase

import (
	"testing"

	"asenare/backend/internal/repository"
)

func TestAuthUseCase_RegisterAndLogin(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	uc := NewAuthUseCase(repo, "test-secret")

	user, token, err := uc.Register("test@example.com", "password123", "テスト")
	if err != nil {
		t.Fatalf("register error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("unexpected email: %s", user.Email)
	}
	if token == "" {
		t.Fatal("token must not be empty")
	}

	_, loginToken, err := uc.Login("test@example.com", "password123")
	if err != nil {
		t.Fatalf("login error: %v", err)
	}
	if loginToken == "" {
		t.Fatal("login token must not be empty")
	}
}

func TestAuthUseCase_LoginFail(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	uc := NewAuthUseCase(repo, "test-secret")

	_, _, _ = uc.Register("test@example.com", "password123", "テスト")

	_, _, err := uc.Login("test@example.com", "wrong")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}
