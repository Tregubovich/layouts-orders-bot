package processor

import (
	"errors"
	"fmt"
	"layouts-orders-bot/internal/clients/client"
	"layouts-orders-bot/internal/entity"
	"layouts-orders-bot/internal/handlers/answers"
	"layouts-orders-bot/internal/handlers/commands"
	"layouts-orders-bot/internal/port"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

type CommandHandler interface {
	HandleCommand(command string, chatID int, username string) (*entity.Message, error)

	NewOrder(username string, props map[entity.State]string) (*entity.Order, error)
	GetOrders(username string) ([]*entity.Message, error)
}

type AnswerHandler interface {
	HandleAnswer(answer string, chatID int) (*entity.Message, error)

	NewSession(chatID int, username string) (*entity.Message, error)
	FinishSession(chatID int) (map[entity.State]string, error)
}

type Processor struct {
	tg       *client.Client
	offset   int
	commands CommandHandler
	answers  AnswerHandler
}

type Meta struct {
	CallbackID string
	ChatID     int
	Username   string
}

func New(client *client.Client, commands CommandHandler, answers AnswerHandler) *Processor {
	return &Processor{
		tg:       client,
		commands: commands,
		answers:  answers,
	}
}

func (p *Processor) Fetch(limit int) ([]entity.Event, error) {
	updates, err := p.tg.GetUpdates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]entity.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, port.EventFromUpdate(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event entity.Event) error {
	switch event.Type {
	case entity.EventTypeCommand:
		return p.processCommand(event)
	case entity.EventTypeAnswer:
		return p.processAnswer(event)
	default:
		return fmt.Errorf("can't process message: %w", ErrUnknownEventType)
	}
}

func (p *Processor) processCommand(event entity.Event) error {
	meta := event.Meta

	p.answers.FinishSession(meta.ChatID)
	response, err := p.commands.HandleCommand(event.Data, meta.ChatID, meta.Username)
	if err != nil {
		if errors.Is(err, commands.ErrNewSession) {
			response, err = p.answers.NewSession(meta.ChatID, meta.Username)
			if err != nil {
				return fmt.Errorf("can't start session: %w", err)
			}
		} else {
			return fmt.Errorf("can't handle command: %w", err)
		}
	}

	if response != nil {
		err = p.tg.SendMessage(meta.ChatID, response.Text, response.Options)
		if err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
	}

	return nil
}

func (p *Processor) processAnswer(event entity.Event) error {
	meta := event.Meta

	if meta.CallbackID != "" {
		if err := p.tg.AnswerCallbackQuery(meta.CallbackID); err != nil {
			return fmt.Errorf("can't answer callback query: %w", err)
		}
	}

	response, err := p.answers.HandleAnswer(event.Data, meta.ChatID)
	if err != nil {
		err = p.proceedAnswerError(err, meta)
		if err != nil {
			return err
		}
	}

	if response != nil {
		err = p.tg.SendMessage(meta.ChatID, response.Text, response.Options)
		if err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
	}

	return nil
}

func (p *Processor) proceedAnswerError(err error, meta entity.Meta) error {
	if errors.Is(err, answers.ErrWrongOption) {
		_ = p.tg.SendMessage(meta.ChatID, fmt.Sprintf("Некорректный ответ: %s", err.Error()), nil)
		return nil
	} else if errors.Is(err, answers.ErrFinishSession) {
		props, err := p.answers.FinishSession(meta.ChatID)
		if err != nil {
			return fmt.Errorf("can't finish session: %w", err)
		}

		order, err := p.commands.NewOrder(meta.Username, props)
		if err != nil {
			return fmt.Errorf("can't create order: %w", err)
		}

		err = p.tg.SendMessage(meta.ChatID, entity.OrderToString(order), nil)
		if err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
		return nil
	}
	return fmt.Errorf("can't handle answer: %w", err)
}
