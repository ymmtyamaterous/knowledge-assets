package repository

import (
	"sort"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryLessonRepository struct {
	mu      sync.RWMutex
	lessons map[string]domain.Lesson
}

func NewInMemoryLessonRepository() *InMemoryLessonRepository {
	now := time.Now().UTC()
	seed := []domain.Lesson{
		{
			ID: "fp3-s1-l1", SectionID: "fp3-s1",
			Title:   "ライフプランとは",
			Content: "# ライフプランとは\n\nライフプランとは、人生のさまざまなイベントを見据えて、長期的な生活設計を行うことです。\n\n## 主なライフイベント\n- 結婚・出産\n- マイホーム購入\n- 子どもの教育\n- 老後の生活\n\n各イベントに必要な資金を試算し、収入・支出のバランスを取ることが重要です。",
			Order:   1, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "fp3-s1-l2", SectionID: "fp3-s1",
			Title:   "家計の収支とバランスシート",
			Content: "# 家計の収支とバランスシート\n\n家計を把握するには「収支」と「資産・負債」の2つの側面から見ます。\n\n## 収支\n- **収入**: 給与・事業収入・不労所得\n- **支出**: 固定費・変動費\n\n## バランスシート（貸借対照表）\n| 資産 | 負債 |\n|------|------|\n| 預貯金 | 住宅ローン |\n| 有価証券 | カードローン |\n| 不動産 | | \n\n純資産 = 資産合計 − 負債合計",
			Order:   2, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "boki3-s1-l1", SectionID: "boki3-s1",
			Title:   "仕訳の基本ルール",
			Content: "# 仕訳の基本ルール\n\n仕訳とは、取引を借方（左側）と貸方（右側）に分けて記録することです。\n\n## 8つの要素\n| | 借方（増加） | 貸方（増加） |\n|--|--|--|\n| 資産 | ✓ | |\n| 負債 | | ✓ |\n| 資本 | | ✓ |\n| 費用 | ✓ | |\n| 収益 | | ✓ |\n\n## 例題\n商品を10,000円で仕入れ、現金で支払った。\n```\n（借）仕入 10,000 / （貸）現金 10,000\n```",
			Order:   1, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "asset3-s1-l1", SectionID: "asset3-s1",
			Title:   "株式投資の基本",
			Content: "# 株式投資の基本\n\n株式とは、企業が資金調達のために発行する証券です。\n\n## 株式の特徴\n- **配当金**: 企業の利益を株主に分配\n- **株主優待**: 企業独自の特典\n- **値上がり益（キャピタルゲイン）**: 株価上昇による利益\n\n## リスク\n株価は市場の需給によって変動します。元本割れのリスクがあります。\n\n## 長期・分散投資\n長期投資と分散投資でリスクを軽減できます。",
			Order:   1, CreatedAt: now, UpdatedAt: now,
		},
	}
	m := make(map[string]domain.Lesson, len(seed))
	for _, l := range seed {
		m[l.ID] = l
	}
	return &InMemoryLessonRepository{lessons: m}
}

func (r *InMemoryLessonRepository) ListBySectionID(sectionID string) ([]domain.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []domain.Lesson
	for _, l := range r.lessons {
		if l.SectionID == sectionID {
			list = append(list, l)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})
	return list, nil
}

func (r *InMemoryLessonRepository) FindByID(id string) (domain.Lesson, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, ok := r.lessons[id]
	return l, ok, nil
}
