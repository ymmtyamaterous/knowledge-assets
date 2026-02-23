package repository

import (
	"sort"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryGlossaryRepository struct {
	mu    sync.RWMutex
	terms map[string]domain.GlossaryTerm
}

func NewInMemoryGlossaryRepository() *InMemoryGlossaryRepository {
	now := time.Now().UTC()
	seed := []domain.GlossaryTerm{
		{ID: "g1", Term: "ライフプラン", Reading: "らいふぷらん", Definition: "人生の各段階における目標を設定し、それを実現するための長期的な生活設計。", CreatedAt: now, UpdatedAt: now},
		{ID: "g2", Term: "仕訳", Reading: "しわけ", Definition: "簿記における取引を借方と貸方に分類して記録する作業。", CreatedAt: now, UpdatedAt: now},
		{ID: "g3", Term: "配当金", Reading: "はいとうきん", Definition: "企業が利益の一部を株主に分配するお金。", CreatedAt: now, UpdatedAt: now},
		{ID: "g4", Term: "投資信託", Reading: "とうししんたく", Definition: "多くの投資家から集めた資金をまとめて運用する金融商品。ファンドとも呼ばれる。", CreatedAt: now, UpdatedAt: now},
		{ID: "g5", Term: "純資産", Reading: "じゅんしさん", Definition: "資産合計から負債合計を差し引いた正味の資産。自己資本とも呼ばれる。", CreatedAt: now, UpdatedAt: now},
	}
	m := make(map[string]domain.GlossaryTerm, len(seed))
	for _, t := range seed {
		m[t.ID] = t
	}
	return &InMemoryGlossaryRepository{terms: m}
}

func (r *InMemoryGlossaryRepository) List() ([]domain.GlossaryTerm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.GlossaryTerm, 0, len(r.terms))
	for _, t := range r.terms {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Reading < list[j].Reading
	})
	return list, nil
}

func (r *InMemoryGlossaryRepository) FindByID(id string) (domain.GlossaryTerm, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.terms[id]
	return t, ok, nil
}
