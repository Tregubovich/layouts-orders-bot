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

		SetMessageID(chatID int, messageID int) error
		CurMessageID(chatID int) (int, error)

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

func (h *Handler) HandleAnswer(answer string, meta entity.Meta) (*entity.Message, error) {
	idx, err := h.repo.CurQuestion(meta.ChatID)
	if err != nil {
		return nil, &NoSessionError{noSessionMsg}
	}

	if meta.MessageID != 0 {
		actual, err := h.repo.CurMessageID(meta.ChatID)
		if err != nil || actual != meta.MessageID {
			return nil, &NoSessionError{noSessionMsg}
		}
	}

	question := Questions[idx]
	value, err := validateOptions(question, answer)
	if err != nil {
		return nil, err
	}

	log.Printf("got answer from '%s: %s", meta.Username, value)

	if question.State == entity.StateAccept {
		if value == AcceptMessage {
			return nil, ErrFinishSession
		}
		return &entity.Message{Text: cancelOrderMsg}, nil
	}

	err = h.repo.SaveAnswer(meta.ChatID, question.State, value)
	if err != nil {
		return nil, fmt.Errorf("can't save answer: %w", err)
	}

	err = h.repo.NextQuestion(meta.ChatID)
	if err != nil {
		return nil, fmt.Errorf("can't set state: %w", err)
	}
	nextQuestion := Questions[idx+1]

	if nextQuestion.State == entity.StateAccept {
		return h.createAcceptQuestion(meta.ChatID)
	}

	return &entity.Message{Text: nextQuestion.Text, Options: FromStringToOptions(nextQuestion.Options)}, nil
}

func (h *Handler) SaveMessageID(messageID int, meta entity.Meta) error {
	return h.repo.SetMessageID(meta.ChatID, messageID)
}

func (h *Handler) GetMessageID(meta entity.Meta) (int, error) {
	return h.repo.CurMessageID(meta.ChatID)
}

func (h *Handler) NewSession(meta entity.Meta) (*entity.Message, error) {
	log.Printf("start session for '%s", meta.Username)

	err := h.repo.StartSession(meta.ChatID)
	if err != nil {
		return nil, fmt.Errorf("can't start session: %w", err)
	}

	firstQuestion := Questions[0]
	return &entity.Message{Text: firstQuestion.Text, Options: FromStringToOptions(firstQuestion.Options)}, nil
}

func (h *Handler) FinishSession(meta entity.Meta) (map[entity.State]string, error) {
	defer h.repo.FinishSession(meta.ChatID)

	props, err := h.repo.GetAnswers(meta.ChatID)
	if err != nil {
		return nil, fmt.Errorf("can't get answers: %w", err)
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

	order := &entity.Order{
		Properties: props,
	}

	minCost, maxCost := h.calculator.Calculate(props)
	acceptQuestion := Questions[len(Questions)-1]
	return &entity.Message{Text: fmt.Sprintf(acceptQuestion.Text, entity.PropertiesToString(order), minCost, maxCost), Options: FromStringToOptions(acceptQuestion.Options)}, nil
}
