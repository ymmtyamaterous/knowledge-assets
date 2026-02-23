package handler

import (
	"errors"
	"net/http"

	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type CourseHandler struct {
	uc *usecase.CourseUseCase
}

func NewCourseHandler(uc *usecase.CourseUseCase) *CourseHandler {
	return &CourseHandler{uc: uc}
}

func (h *CourseHandler) List(w http.ResponseWriter, _ *http.Request) {
	courses, err := h.uc.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{"courses": courses})
}

func (h *CourseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	course, err := h.uc.Get(id)
	if err != nil {
		if errors.Is(err, usecase.ErrCourseNotFound) {
			WriteError(w, http.StatusNotFound, "course not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, course)
}
