package usecase

import (
	"testing"

	"asenare/backend/internal/repository"
)

func TestGlossaryUseCase_ListWithTag(t *testing.T) {
	repo := repository.NewInMemoryGlossaryRepository()
	uc := NewGlossaryUseCase(repo)

	tags, err := uc.ListTags()
	if err != nil {
		t.Fatalf("list tags error: %v", err)
	}
	if len(tags) == 0 {
		t.Fatal("tags should not be empty")
	}

	terms, err := uc.List("tag-boki")
	if err != nil {
		t.Fatalf("list with tag error: %v", err)
	}
	if len(terms) == 0 {
		t.Fatal("terms should not be empty")
	}
	for _, term := range terms {
		if len(term.Tags) == 0 {
			t.Fatal("term should have tags")
		}
	}
}
