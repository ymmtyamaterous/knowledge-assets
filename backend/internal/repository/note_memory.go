package repository

import (
	"fmt"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type MemoryNoteRepository struct {
	mu    sync.RWMutex
	notes map[string]domain.UserNote // key: userID + ":" + lessonID
}

func NewMemoryNoteRepository() *MemoryNoteRepository {
	return &MemoryNoteRepository{notes: make(map[string]domain.UserNote)}
}

func noteKey(userID, lessonID string) string {
	return userID + ":" + lessonID
}

func (r *MemoryNoteRepository) FindByUserAndLesson(userID, lessonID string) (domain.UserNote, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.notes[noteKey(userID, lessonID)]
	return n, ok, nil
}

func (r *MemoryNoteRepository) Upsert(note domain.UserNote) (domain.UserNote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	key := noteKey(note.UserID, note.LessonID)
	if existing, ok := r.notes[key]; ok {
		existing.Content = note.Content
		existing.UpdatedAt = now
		r.notes[key] = existing
		return existing, nil
	}
	if note.ID == "" {
		note.ID = fmt.Sprintf("note_%d", now.UnixNano())
	}
	note.CreatedAt = now
	note.UpdatedAt = now
	r.notes[key] = note
	return note, nil
}

func (r *MemoryNoteRepository) ListByUserID(userID string) ([]domain.UserNote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.UserNote
	for _, n := range r.notes {
		if n.UserID == userID {
			result = append(result, n)
		}
	}
	return result, nil
}
