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
		// FP3級
		{ID: "fp3-s1", CourseID: "fp3", Title: "第1章: ライフプランと資金計画", Description: "人生設計とお金の基礎", Order: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "fp3-s2", CourseID: "fp3", Title: "第2章: 保険の基礎", Description: "生命保険・損害保険・社会保険の仕組み", Order: 2, CreatedAt: now, UpdatedAt: now},
		{ID: "fp3-s3", CourseID: "fp3", Title: "第3章: 金融資産運用", Description: "預金・株式・債券・投資信託の基礎知識", Order: 3, CreatedAt: now, UpdatedAt: now},
		{ID: "fp3-s4", CourseID: "fp3", Title: "第4章: タックスプランニング", Description: "所得税・控除・NISA・iDeCoの活用法", Order: 4, CreatedAt: now, UpdatedAt: now},
		{ID: "fp3-s5", CourseID: "fp3", Title: "第5章: 不動産と相続", Description: "不動産取引・相続・贈与の基礎知識", Order: 5, CreatedAt: now, UpdatedAt: now},
		// 簿記3級
		{ID: "boki3-s1", CourseID: "boki3", Title: "第1章: 簿記の基本", Description: "仕訳・勘定科目の基礎", Order: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "boki3-s2", CourseID: "boki3", Title: "第2章: 主要な勘定科目", Description: "現金・売掛金・固定資産などの処理方法", Order: 2, CreatedAt: now, UpdatedAt: now},
		{ID: "boki3-s3", CourseID: "boki3", Title: "第3章: 決算と財務諸表", Description: "試算表・損益計算書・貸借対照表の読み方", Order: 3, CreatedAt: now, UpdatedAt: now},
		// 資産運用検定3級
		{ID: "asset3-s1", CourseID: "asset3", Title: "第1章: 投資の基本", Description: "株式・債券・投資信託の概要", Order: 1, CreatedAt: now, UpdatedAt: now},
		{ID: "asset3-s2", CourseID: "asset3", Title: "第2章: 債券と投資信託", Description: "債券の仕組みと投資信託・ETFの選び方", Order: 2, CreatedAt: now, UpdatedAt: now},
		{ID: "asset3-s3", CourseID: "asset3", Title: "第3章: ポートフォリオとリスク管理", Description: "分散投資・長期投資の考え方", Order: 3, CreatedAt: now, UpdatedAt: now},
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
