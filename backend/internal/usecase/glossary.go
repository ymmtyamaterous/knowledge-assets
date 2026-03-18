package usecase

import (
	"errors"
	"time"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

type GlossaryUseCase struct {
	glossary repository.GlossaryRepository
}

var ErrGlossaryTermNotFound = errors.New("glossary term not found")

func NewGlossaryUseCase(glossary repository.GlossaryRepository) *GlossaryUseCase {
	return &GlossaryUseCase{glossary: glossary}
}

func (uc *GlossaryUseCase) List(tagID string) ([]domain.GlossaryTerm, error) {
	return uc.glossary.List(tagID)
}

func (uc *GlossaryUseCase) Get(id string) (domain.GlossaryTerm, error) {
	t, ok, err := uc.glossary.FindByID(id)
	if err != nil {
		return domain.GlossaryTerm{}, err
	}
	if !ok {
		return domain.GlossaryTerm{}, ErrGlossaryTermNotFound
	}
	return t, nil
}

func (uc *GlossaryUseCase) ListTags() ([]domain.GlossaryTag, error) {
	return uc.glossary.ListTags()
}

// GetDailyTerm は日付ベースで決定論的に1つの用語を返します。
func (uc *GlossaryUseCase) GetDailyTerm() (domain.GlossaryTerm, error) {
	terms, err := uc.glossary.List("")
	if err != nil {
		return domain.GlossaryTerm{}, err
	}
	if len(terms) == 0 {
		return domain.GlossaryTerm{}, ErrGlossaryTermNotFound
	}
	// 今日の日付のシリアル番号（日単位）を用語数で割った余りでインデックスを決定
	today := time.Now().Truncate(24 * time.Hour)
	daySerial := today.Unix() / 86400
	idx := int(daySerial) % len(terms)
	return terms[idx], nil
}
