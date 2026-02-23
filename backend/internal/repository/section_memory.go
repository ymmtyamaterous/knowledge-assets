package repository

import (
	"sort"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemorySectionRepository struct {
	mu       sync.RWMutex
	sections map[string]domain.Section
}

func NewInMemorySectionRepository() *InMemorySectionRepository {
	now := time.Now().UTC()
	seed := []domain.Section{
		{ID: "fp3-s1", CourseID: "fp3", Title: "第1章: ライフプランと資金計画", Description: "人生設計とお金の基礎", Order: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "fp3-s2", CourseID: "fp3", Title: "第2章: 保険の基礎", Description: "生命保険・損害保険の仕組み", Order: 2, CreatedAt: now, UpdatedAt: now},
		{ID: "boki3-s1", CourseID: "boki3", Title: "第1章: 簿記の基本", Description: "仕訳・勘定科目の基礎", Order: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "asset3-s1", CourseID: "asset3", Title: "第1章: 投資の基本", Description: "株式・債券・投資信託の概要", Order: 1, CreatedAt: now, UpdatedAt: now},
	}
	m := make(map[string]domain.Section, len(seed))
	for _, s := range seed {
		m[s.ID] = s
	}
	return &InMemorySectionRepository{sections: m}
}

func (r *InMemorySectionRepository) ListByCourseID(courseID string) ([]domain.Section, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []domain.Section
	for _, s := range r.sections {
		if s.CourseID == courseID {
			list = append(list, s)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})
	return list, nil
}

func (r *InMemorySectionRepository) FindByID(id string) (domain.Section, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sections[id]
	return s, ok, nil
}
