package handler

import (
	"errors"
	"net/http"

	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type SectionHandler struct {
	uc *usecase.SectionUseCase
}

func NewSectionHandler(uc *usecase.SectionUseCase) *SectionHandler {
	return &SectionHandler{uc: uc}
}

func (h *SectionHandler) ListByCourse(w http.ResponseWriter, r *http.Request) {
	courseID := chi.URLParam(r, "courseID")
	sections, err := h.uc.ListByCourse(courseID)
	if err != nil {
		if errors.Is(err, usecase.ErrCourseNotFound) {
			WriteError(w, http.StatusNotFound, "course not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"sections": sections})
}
