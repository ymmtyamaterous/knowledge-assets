package repository

import (
	"sort"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryCourseRepository struct {
	mu      sync.RWMutex
	courses map[string]domain.Course
}

func NewInMemoryCourseRepository() *InMemoryCourseRepository {
	now := time.Now().UTC()
	seed := []domain.Course{
		{ID: "fp3", Title: "FP3級", Description: "家計・保険・税金の基礎を学ぶ", Difficulty: "beginner", EstimatedHour: 30, Order: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "boki3", Title: "簿記3級", Description: "会計・仕訳の基本を学ぶ", Difficulty: "beginner", EstimatedHour: 35, Order: 2, CreatedAt: now, UpdatedAt: now},
		{ID: "asset3", Title: "資産運用検定3級", Description: "投資と資産形成の入門を学ぶ", Difficulty: "beginner", EstimatedHour: 25, Order: 3, CreatedAt: now, UpdatedAt: now},
	}

	m := make(map[string]domain.Course, len(seed))
	for _, c := range seed {
		m[c.ID] = c
	}

	return &InMemoryCourseRepository{courses: m}
}

func (r *InMemoryCourseRepository) List() ([]domain.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Course, 0, len(r.courses))
	for _, c := range r.courses {
		list = append(list, c)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})

	return list, nil
}

func (r *InMemoryCourseRepository) FindByID(id string) (domain.Course, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.courses[id]
	return c, ok, nil
}
