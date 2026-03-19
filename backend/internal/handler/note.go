package handler

import (
	"encoding/json"
	"net/http"

	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	uc *usecase.NoteUseCase
}

func NewNoteHandler(uc *usecase.NoteUseCase) *NoteHandler {
	return &NoteHandler{uc: uc}
}

func (h *NoteHandler) GetByLesson(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonID := chi.URLParam(r, "lessonId")
	note, ok, err := h.uc.GetNote(userID, lessonID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if !ok {
		WriteJSON(w, http.StatusOK, map[string]any{"note": nil})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"note": note})
}

func (h *NoteHandler) Save(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lessonID := chi.URLParam(r, "lessonId")

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := h.uc.SaveNote(userID, lessonID, req.Content)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"note": note})
}

func (h *NoteHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notes, err := h.uc.ListNotesWithLesson(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"notes": notes})
}
