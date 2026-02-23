package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type QuizHandler struct {
	uc *usecase.QuizUseCase
}

func NewQuizHandler(uc *usecase.QuizUseCase) *QuizHandler {
	return &QuizHandler{uc: uc}
}

func (h *QuizHandler) GetByLesson(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonId")
	quiz, err := h.uc.FindByLessonID(lessonID)
	if err != nil {
		if errors.Is(err, usecase.ErrLessonNotFound) {
			WriteError(w, http.StatusNotFound, "lesson not found")
			return
		}
		if errors.Is(err, usecase.ErrQuizNotFound) {
			WriteError(w, http.StatusNotFound, "quiz not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, quiz)
}

func (h *QuizHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	detail, err := h.uc.Get(id)
	if err != nil {
		if errors.Is(err, usecase.ErrQuizNotFound) {
			WriteError(w, http.StatusNotFound, "quiz not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	WriteJSON(w, http.StatusOK, detail)
}

type submitQuizRequest struct {
	Answers []usecase.QuizAnswer `json:"answers"`
}

func (h *QuizHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	quizID := chi.URLParam(r, "id")

	var req submitQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request")
		return
	}

	result, err := h.uc.Submit(userID, quizID, req.Answers)
	if err != nil {
		if errors.Is(err, usecase.ErrQuizNotFound) {
			WriteError(w, http.StatusNotFound, "quiz not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	WriteJSON(w, http.StatusOK, result)
}

func (h *QuizHandler) ListMyResults(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(userIDContextKey).(string)
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	list, err := h.uc.ListResults(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if list == nil {
		list = []domain.UserQuizResult{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"results": list})
}
