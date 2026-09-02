package answers

import (
	"fmt"
	"layouts-orders-bot/internal/entity"
	"log"
	"slices"
	"strconv"
	"strings"
)

const (
	cancelOrderMsg = "Заказ отменён"
	noSessionMsg   = "Сессия не валидна"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler_mocks.go -package=mocks
type (
	Repository interface {
		CurQuestion(chatID int) (int, error)
		NextQuestion(chatID int) error

		SaveAnswer(chatID int, key entity.State, value string) error
		GetAnswers(chatID int) (map[entity.State]string, error)

		StartSession(chatID int) error
		FinishSession(chatID int) error
	}

	Calculator interface {
		Calculate(props map[entity.State]string) (int, int)
	}
)

type Handler struct {
	repo       Repository
	calculator Calculator
}

func NewHandler(storage Repository, calculator Calculator) *Handler {
	return &Handler{
		repo:       storage,
		calculator: calculator,
	}
}

func (h *Handler) HandleAnswer(answer string, chatID int, username string) (*entity.Message, error) {
	idx, err := h.repo.CurQuestion(chatID)
	if err != nil {
		return nil, &NoSessionError{noSessionMsg}
	}

	question := Questions[idx]
	value, err := validateOptions(question, answer)
	if err != nil {
		return nil, err
	}

	log.Printf("got answer from %s: %s", username, value)

	if question.State == entity.StateAccept {
		if value == AcceptMessage {
			return nil, ErrFinishSession
		}
		return &entity.Message{Text: cancelOrderMsg}, nil
	}

	err = h.repo.SaveAnswer(chatID, question.State, value)
	if err != nil {
		return nil, fmt.Errorf("can't save answer: %w", err)
	}

	err = h.repo.NextQuestion(chatID)
	if err != nil {
		return nil, fmt.Errorf("can't set state: %w", err)
	}
	nextQuestion := Questions[idx+1]

	if nextQuestion.State == entity.StateAccept {
		return h.createAcceptQuestion(chatID)
	}

	return &entity.Message{Text: nextQuestion.Text, Options: FromStringToOptions(nextQuestion.Options)}, nil
}

func (h *Handler) NewSession(chatID int, username string) (*entity.Message, error) {
	log.Printf("start session for %s", username)

	err := h.repo.StartSession(chatID)
	if err != nil {
		return nil, fmt.Errorf("can't start session: %w", err)
	}

	firstQuestion := Questions[0]
	return &entity.Message{Text: firstQuestion.Text, Options: FromStringToOptions(firstQuestion.Options)}, nil
}

func (h *Handler) FinishSession(chatID int) (map[entity.State]string, error) {
	props, err := h.repo.GetAnswers(chatID)
	if err != nil {
		return nil, fmt.Errorf("can't get answers: %w", err)
	}

	err = h.repo.FinishSession(chatID)
	if err != nil {
		return nil, fmt.Errorf("can't finish session: %w", err)
	}

	log.Printf("got answers: %+v", props)

	return props, nil
}

func validateOptions(question *Question, answer string) (string, error) {
	options := question.Options
	num, err := strconv.Atoi(answer)
	if err == nil {
		if num >= 1 && num <= len(options) {
			return options[num-1], nil
		}
	}
	if !slices.Contains(options, answer) {
		if question.SpecialValidate == nil {
			return "", &WrongOptionError{fmt.Sprintf("должен быть один из (%s)", strings.Join(options, ", "))}
		}
		return question.SpecialValidate(answer)
	}
	return answer, nil
}

func (h *Handler) createAcceptQuestion(chatID int) (*entity.Message, error) {
	props, err := h.repo.GetAnswers(chatID)
	if err != nil {
		return nil, fmt.Errorf("can't get answers: %w", err)
	}

	minCost, maxCost := h.calculator.Calculate(props)
	acceptQuestion := Questions[len(Questions)-1]
	return &entity.Message{Text: fmt.Sprintf(acceptQuestion.Text, minCost, maxCost), Options: FromStringToOptions(acceptQuestion.Options)}, nil
}
