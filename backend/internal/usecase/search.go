package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var ErrSearchQueryTooShort = errors.New("search query must be at least 2 characters")

type SearchUseCase struct {
	repo repository.SearchRepository
}

func NewSearchUseCase(repo repository.SearchRepository) *SearchUseCase {
	return &SearchUseCase{repo: repo}
}

func (uc *SearchUseCase) Search(query string) (domain.SearchResult, error) {
	if len([]rune(query)) < 2 {
		return domain.SearchResult{}, ErrSearchQueryTooShort
	}

	lessons, err := uc.repo.SearchLessons(query)
	if err != nil {
		return domain.SearchResult{}, err
	}
	terms, err := uc.repo.SearchTerms(query)
	if err != nil {
		return domain.SearchResult{}, err
	}

	if lessons == nil {
		lessons = []domain.SearchLesson{}
	}
	if terms == nil {
		terms = []domain.SearchTerm{}
	}

	return domain.SearchResult{Lessons: lessons, Terms: terms}, nil
}
