package repository

import (
	"fmt"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryBadgeRepository struct {
	mu         sync.RWMutex
	badges     map[string]domain.Badge     // key: id
	userBadges map[string]domain.UserBadge // key: userID:badgeID
}

func NewInMemoryBadgeRepository() *InMemoryBadgeRepository {
	return &InMemoryBadgeRepository{
		badges:     make(map[string]domain.Badge),
		userBadges: make(map[string]domain.UserBadge),
	}
}

// SeedBadge はテスト用にバッジをセットするヘルパーです。
func (r *InMemoryBadgeRepository) SeedBadge(b domain.Badge) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.badges[b.ID] = b
}

func (r *InMemoryBadgeRepository) FindByCondition(conditionType, conditionID string) (domain.Badge, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.badges {
		if b.ConditionType == conditionType && b.ConditionID == conditionID {
			return b, true, nil
		}
	}
	return domain.Badge{}, false, nil
}

func (r *InMemoryBadgeRepository) CreateUserBadge(userID, badgeID string) (domain.UserBadge, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := userID + ":" + badgeID
	if ub, ok := r.userBadges[key]; ok {
		return ub, nil // 冪等
	}
	badge, ok := r.badges[badgeID]
	if !ok {
		return domain.UserBadge{}, fmt.Errorf("badge not found: %s", badgeID)
	}
	ub := domain.UserBadge{
		ID:       fmt.Sprintf("ub-%d", time.Now().UnixNano()),
		UserID:   userID,
		Badge:    badge,
		EarnedAt: time.Now().UTC(),
	}
	r.userBadges[key] = ub
	return ub, nil
}

func (r *InMemoryBadgeRepository) ListByUserID(userID string) ([]domain.UserBadge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.UserBadge
	for key, ub := range r.userBadges {
		if len(key) > len(userID) && key[:len(userID)] == userID && key[len(userID)] == ':' {
			result = append(result, ub)
		}
	}
	return result, nil
}

func (r *InMemoryBadgeRepository) ExistsUserBadge(userID, badgeID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.userBadges[userID+":"+badgeID]
	return ok, nil
}
