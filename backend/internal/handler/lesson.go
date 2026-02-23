package handler

import (
	"errors"
	"net/http"

	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type LessonHandler struct {
	uc *usecase.LessonUseCase
}

func NewLessonHandler(uc *usecase.LessonUseCase) *LessonHandler {
	return &LessonHandler{uc: uc}
}

func (h *LessonHandler) ListBySection(w http.ResponseWriter, r *http.Request) {
	sectionID := chi.URLParam(r, "sectionID")
	lessons, err := h.uc.ListBySection(sectionID)
	if err != nil {
		if errors.Is(err, usecase.ErrSectionNotFound) {
			WriteError(w, http.StatusNotFound, "section not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"lessons": lessons})
}

func (h *LessonHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lesson, err := h.uc.Get(id)
	if err != nil {
		if errors.Is(err, usecase.ErrLessonNotFound) {
			WriteError(w, http.StatusNotFound, "lesson not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, lesson)
}
