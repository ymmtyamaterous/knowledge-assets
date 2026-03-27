package usecase

import (
	"errors"
	"strings"
	"time"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

type AuthUseCase struct {
	users     repository.UserRepository
	jwtSecret string
}

func NewAuthUseCase(users repository.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{users: users, jwtSecret: jwtSecret}
}

func (uc *AuthUseCase) Register(email, password, username string) (domain.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)
	if email == "" || password == "" {
		return domain.User{}, "", ErrInvalidCredentials
	}
	if username == "" {
		username = "ユーザー"
	}

	if _, exists, err := uc.users.FindByEmail(email); err != nil {
		return domain.User{}, "", err
	} else if exists {
		return domain.User{}, "", ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, "", err
	}

	user, err := uc.users.Create(domain.User{
		Email:        email,
		PasswordHash: string(hash),
		Username:     username,
		Role:         domain.RoleUser,
	})
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return domain.User{}, "", ErrEmailExists
		}
		return domain.User{}, "", err
	}

	token, err := uc.issueToken(user)
	if err != nil {
		return domain.User{}, "", err
	}

	return user, token, nil
}

func (uc *AuthUseCase) Login(email, password string) (domain.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, exists, err := uc.users.FindByEmail(email)
	if err != nil {
		return domain.User{}, "", err
	}
	if !exists {
		return domain.User{}, "", ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return domain.User{}, "", ErrInvalidCredentials
	}

	token, err := uc.issueToken(user)
	if err != nil {
		return domain.User{}, "", err
	}

	return user, token, nil
}

func (uc *AuthUseCase) ChangePassword(userID, currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return ErrInvalidCredentials
	}

	user, ok, err := uc.users.FindByID(userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)) != nil {
		return ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	_, err = uc.users.Update(user)
	return err
}

func (uc *AuthUseCase) issueToken(user domain.User) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": string(user.Role),
		"iat":  now.Unix(),
		"exp":  now.Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
