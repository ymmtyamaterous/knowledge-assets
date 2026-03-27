package handler

import (
	"errors"
	"net/http"

	"asenare/backend/internal/usecase"
)

type SearchHandler struct {
	uc *usecase.SearchUseCase
}

func NewSearchHandler(uc *usecase.SearchUseCase) *SearchHandler {
	return &SearchHandler{uc: uc}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := h.uc.Search(q)
	if err != nil {
		if errors.Is(err, usecase.ErrSearchQueryTooShort) {
			WriteError(w, http.StatusBadRequest, "search query must be at least 2 characters")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, result)
}
