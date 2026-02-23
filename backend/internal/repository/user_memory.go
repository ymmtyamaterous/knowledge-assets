package repository

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryUserRepository struct {
	mu      sync.RWMutex
	byID    map[string]domain.User
	byEmail map[string]string
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byID:    map[string]domain.User{},
		byEmail: map[string]string{},
	}
}

func (r *InMemoryUserRepository) Create(user domain.User) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	email := strings.ToLower(strings.TrimSpace(user.Email))
	if email == "" {
		return domain.User{}, errors.New("email is required")
	}
	if _, exists := r.byEmail[email]; exists {
		return domain.User{}, errors.New("email already exists")
	}

	now := time.Now().UTC()
	if user.ID == "" {
		user.ID = fmt.Sprintf("u_%d", now.UnixNano())
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now
	user.Email = email

	r.byID[user.ID] = user
	r.byEmail[email] = user.ID

	return user, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (domain.User, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return domain.User{}, false, nil
	}
	u, ok := r.byID[id]
	return u, ok, nil
}

func (r *InMemoryUserRepository) FindByID(id string) (domain.User, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.byID[id]
	return u, ok, nil
}

func (r *InMemoryUserRepository) Update(user domain.User) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[user.ID]; !ok {
		return domain.User{}, errors.New("user not found")
	}
	user.UpdatedAt = time.Now().UTC()
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user.ID
	return user, nil
}
