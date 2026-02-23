package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrGlossaryTermNotFound = errors.New("glossary term not found")

type GlossaryUseCase struct {
	glossary repository.GlossaryRepository
}

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
