package usecase

import (
	"errors"

	"asenare/backend/internal/domain"
	"asenare/backend/internal/repository"
)

var (
	ErrQuizNotFound = errors.New("quiz not found")
)

type QuizUseCase struct {
	quizzes repository.QuizRepository
	lessons repository.LessonRepository
}

func NewQuizUseCase(quizzes repository.QuizRepository, lessons repository.LessonRepository) *QuizUseCase {
	return &QuizUseCase{quizzes: quizzes, lessons: lessons}
}

type QuizDetail struct {
	Quiz      domain.Quiz           `json:"quiz"`
	Questions []domain.QuizQuestion `json:"questions"`
}

type QuizAnswer struct {
	QuestionID string `json:"questionId"`
	ChoiceID   string `json:"choiceId"`
}

type SubmitQuizResult struct {
	Result  domain.UserQuizResult        `json:"result"`
	Correct map[string]bool              `json:"correct"`
	Answers map[string]string            `json:"answers"`
	Details map[string]domain.QuizChoice `json:"details"`
}

func (uc *QuizUseCase) FindByLessonID(lessonID string) (domain.Quiz, error) {
	if _, ok, err := uc.lessons.FindByID(lessonID); err != nil {
		return domain.Quiz{}, err
	} else if !ok {
		return domain.Quiz{}, ErrLessonNotFound
	}

	quiz, ok, err := uc.quizzes.FindByLessonID(lessonID)
	if err != nil {
		return domain.Quiz{}, err
	}
	if !ok {
		return domain.Quiz{}, ErrQuizNotFound
	}
	return quiz, nil
}

func (uc *QuizUseCase) Get(id string) (QuizDetail, error) {
	quiz, ok, err := uc.quizzes.FindByID(id)
	if err != nil {
		return QuizDetail{}, err
	}
	if !ok {
		return QuizDetail{}, ErrQuizNotFound
	}

	questions, err := uc.quizzes.ListQuestions(id)
	if err != nil {
		return QuizDetail{}, err
	}

	return QuizDetail{Quiz: quiz, Questions: questions}, nil
}

func (uc *QuizUseCase) Submit(userID, quizID string, answers []QuizAnswer) (SubmitQuizResult, error) {
	_, ok, err := uc.quizzes.FindByID(quizID)
	if err != nil {
		return SubmitQuizResult{}, err
	}
	if !ok {
		return SubmitQuizResult{}, ErrQuizNotFound
	}

	questions, err := uc.quizzes.ListQuestions(quizID)
	if err != nil {
		return SubmitQuizResult{}, err
	}

	answerByQuestion := make(map[string]string, len(answers))
	for _, answer := range answers {
		answerByQuestion[answer.QuestionID] = answer.ChoiceID
	}

	correct := 0
	correctByQuestion := make(map[string]bool, len(questions))
	details := make(map[string]domain.QuizChoice)
	for _, question := range questions {
		selectedChoiceID := answerByQuestion[question.ID]
		for _, choice := range question.Choices {
			if choice.ID == selectedChoiceID {
				details[question.ID] = choice
				if choice.IsCorrect {
					correct++
					correctByQuestion[question.ID] = true
				} else {
					correctByQuestion[question.ID] = false
				}
				break
			}
		}
	}

	result, err := uc.quizzes.CreateResult(domain.UserQuizResult{
		UserID: userID,
		QuizID: quizID,
		Score:  correct,
		Total:  len(questions),
	})
	if err != nil {
		return SubmitQuizResult{}, err
	}

	return SubmitQuizResult{
		Result:  result,
		Correct: correctByQuestion,
		Answers: answerByQuestion,
		Details: details,
	}, nil
}

func (uc *QuizUseCase) ListResults(userID string) ([]domain.UserQuizResult, error) {
	return uc.quizzes.ListResultsByUserID(userID)
}
