package usecase

import (
	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

type NoteUseCase struct {
	notes   repository.NoteRepository
	lessons repository.LessonRepository
}

func NewNoteUseCase(notes repository.NoteRepository, lessons repository.LessonRepository) *NoteUseCase {
	return &NoteUseCase{notes: notes, lessons: lessons}
}

func (uc *NoteUseCase) GetNote(userID, lessonID string) (domain.UserNote, bool, error) {
	return uc.notes.FindByUserAndLesson(userID, lessonID)
}

func (uc *NoteUseCase) SaveNote(userID, lessonID, content string) (domain.UserNote, error) {
	if _, ok, err := uc.lessons.FindByID(lessonID); err != nil {
		return domain.UserNote{}, err
	} else if !ok {
		return domain.UserNote{}, ErrLessonNotFound
	}

	return uc.notes.Upsert(domain.UserNote{
		UserID:   userID,
		LessonID: lessonID,
		Content:  content,
	})
}

func (uc *NoteUseCase) ListNotes(userID string) ([]domain.UserNote, error) {
	notes, err := uc.notes.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	if notes == nil {
		return []domain.UserNote{}, nil
	}
	return notes, nil
}
