package repository

import (
	"fmt"
	"sync"
	"time"

	"asenare/backend/internal/domain"
)

type InMemoryProgressRepository struct {
	mu       sync.RWMutex
	progress map[string]domain.UserLessonProgress // key: userID+":"+lessonID
}

func NewInMemoryProgressRepository() *InMemoryProgressRepository {
	return &InMemoryProgressRepository{
		progress: map[string]domain.UserLessonProgress{},
	}
}

func progressKey(userID, lessonID string) string {
	return userID + ":" + lessonID
}

func (r *InMemoryProgressRepository) FindByUserAndLesson(userID, lessonID string) (domain.UserLessonProgress, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.progress[progressKey(userID, lessonID)]
	return p, ok, nil
}

func (r *InMemoryProgressRepository) Create(p domain.UserLessonProgress) (domain.UserLessonProgress, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	if p.ID == "" {
		p.ID = fmt.Sprintf("prog_%d", now.UnixNano())
	}
	if p.CompletedAt.IsZero() {
		p.CompletedAt = now
	}

	r.progress[progressKey(p.UserID, p.LessonID)] = p
	return p, nil
}

func (r *InMemoryProgressRepository) DeleteByUserAndLesson(userID, lessonID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.progress, progressKey(userID, lessonID))
	return nil
}

func (r *InMemoryProgressRepository) ListByUserID(userID string) ([]domain.UserLessonProgress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []domain.UserLessonProgress
	for _, p := range r.progress {
		if p.UserID == userID {
			list = append(list, p)
		}
	}
	return list, nil
}

func (r *InMemoryProgressRepository) CountByUserAndCourse(userID, _ string, lessonIDs []string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	set := make(map[string]struct{}, len(lessonIDs))
	for _, id := range lessonIDs {
		set[id] = struct{}{}
	}

	count := 0
	for _, p := range r.progress {
		if p.UserID == userID {
			if _, ok := set[p.LessonID]; ok {
				count++
			}
		}
	}
	return count, nil
}
