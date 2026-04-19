package handler

import (
	"errors"
	"fmt"
	"net/http"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type ProgressHandler struct {
	uc *usecase.ProgressUseCase
}

func NewProgressHandler(uc *usecase.ProgressUseCase) *ProgressHandler {
	return &ProgressHandler{uc: uc}
}

func (h *ProgressHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonID := chi.URLParam(r, "id")
	result, err := h.uc.CompleteLesson(userID, lessonID)
	if err != nil {
		if errors.Is(err, usecase.ErrLessonNotFound) {
			WriteError(w, http.StatusNotFound, "lesson not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, result)
}

func (h *ProgressHandler) UncompleteLesson(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonID := chi.URLParam(r, "id")
	if err := h.uc.UncompleteLesson(userID, lessonID); err != nil {
		if errors.Is(err, usecase.ErrLessonNotFound) {
			WriteError(w, http.StatusNotFound, "lesson not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *ProgressHandler) GetMyProgress(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	list, err := h.uc.GetUserProgress(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if list == nil {
		list = []domain.UserLessonProgress{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"progress": list})
}

func (h *ProgressHandler) GetMyCourseProgress(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	list, err := h.uc.GetCourseProgress(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if list == nil {
		list = []domain.CourseProgress{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"courseProgress": list})
}

func (h *ProgressHandler) GetMyStreak(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	streak, err := h.uc.GetStreak(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, streak)
}

func (h *ProgressHandler) GetMyStats(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	stats, err := h.uc.GetStats(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, stats)
}

func (h *ProgressHandler) GetMyCalendar(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var year int
	if y := r.URL.Query().Get("year"); y != "" {
		if _, err := fmt.Sscanf(y, "%d", &year); err != nil {
			year = 0
		}
	}

	cal, err := h.uc.GetCalendar(userID, year)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, cal)
}

func (h *ProgressHandler) GetMyBadges(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	badges, err := h.uc.GetMyBadges(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if badges == nil {
		badges = []domain.UserBadge{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"badges": badges})
}
