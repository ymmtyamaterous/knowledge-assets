package handler

import (
	"errors"
	"net/http"

	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type GlossaryHandler struct {
	uc *usecase.GlossaryUseCase
}

func NewGlossaryHandler(uc *usecase.GlossaryUseCase) *GlossaryHandler {
	return &GlossaryHandler{uc: uc}
}

func (h *GlossaryHandler) List(w http.ResponseWriter, _ *http.Request) {
	terms, err := h.uc.List()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"terms": terms})
}

func (h *GlossaryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	term, err := h.uc.Get(id)
	if err != nil {
		if errors.Is(err, usecase.ErrGlossaryTermNotFound) {
			WriteError(w, http.StatusNotFound, "term not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, term)
}
