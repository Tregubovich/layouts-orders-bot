package processor

import (
	"errors"
	"fmt"
	"layouts-orders-bot/internal/clients/telegram"
	"layouts-orders-bot/internal/entity"
	"layouts-orders-bot/internal/handlers/answers"
	"layouts-orders-bot/internal/handlers/commands"
	"layouts-orders-bot/internal/port"
	"log"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

type CommandHandler interface {
	HandleCommand(command string, meta entity.Meta) (*entity.Message, error)

	NewOrder(props map[*entity.State]string, meta entity.Meta) (*entity.Message, error)
	GetOrders(meta entity.Meta, isAdmin bool) ([]*entity.Message, error)
}

type AnswerHandler interface {
	HandleAnswer(answer string, meta entity.Meta) (*entity.Message, error)

	SaveMessageID(messageID int, meta entity.Meta) error
	GetMessageID(meta entity.Meta) (int, error)

	NewSession(meta entity.Meta) (*entity.Message, error)
	FinishSession(meta entity.Meta) (map[*entity.State]string, error)
}

type Processor struct {
	tg       *telegram.Client
	offset   int
	commands CommandHandler
	answers  AnswerHandler
}

func New(client *telegram.Client, commands CommandHandler, answers AnswerHandler) *Processor {
	p := &Processor{
		tg:       client,
		commands: commands,
		answers:  answers,
	}

	err := p.tg.SetCommands(entity.Commands)
	if err != nil {
		log.Printf("telegram: failed to set commands: %v", err)
	}

	return p
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

	id, _ := p.answers.GetMessageID(meta)
	if id != 0 {
		if err := p.tg.RemoveKeyboard(meta.ChatID, id); err != nil {
			log.Printf("can't delete message: %v", err)
		}
	}
	_, _ = p.answers.FinishSession(meta)

	response, err := p.commands.HandleCommand(event.Data, meta)
	if err != nil {
		err = p.processCommandError(err, response, meta)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = p.tg.SendMessage(meta.ChatID, response.Text, response.Options)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (p *Processor) processCommandError(err error, response *entity.Message, meta entity.Meta) error {
	if errors.Is(err, commands.ErrNewSession) {
		response, err = p.answers.NewSession(meta)
		if err != nil {
			return fmt.Errorf("can't start session: %w", err)
		}

		id, err := p.tg.SendMessage(meta.ChatID, response.Text, response.Options)
		if err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}

		return p.answers.SaveMessageID(id, meta)
	} else if errors.Is(err, commands.ErrGetOrders) || errors.Is(err, commands.ErrGetAllOrders) {
		var isAdmin bool
		if errors.Is(err, commands.ErrGetOrders) {
			isAdmin = false
		} else {
			isAdmin = true
		}

		messages, err := p.commands.GetOrders(meta, isAdmin)
		if err != nil {
			return fmt.Errorf("can't get orders: %w", err)
		}

		for _, message := range messages {
			_, err = p.tg.SendMessage(meta.ChatID, message.Text, message.Options)
			if err != nil {
				log.Printf("can't send message: %v", err)
			}
		}

		return nil
	} else {
		return fmt.Errorf("can't handle command: %w", err)
	}
}

func (p *Processor) processAnswer(event entity.Event) error {
	meta := event.Meta

	if meta.CallbackID != "" {
		if err := p.tg.AnswerCallbackQuery(meta.CallbackID); err != nil {
			return fmt.Errorf("can't answer callback query: %w", err)
		}
	}

	response, err := p.answers.HandleAnswer(event.Data, meta)
	if err != nil {
		err = p.proceedAnswerError(err, meta)
		if err != nil {
			return err
		}
	}

	if meta.CallbackID != "" {
		if err := p.tg.RemoveKeyboard(meta.ChatID, meta.MessageID); err != nil {
			log.Printf("can't delete message: %v", err)
		}
	}

	if response != nil {
		p.removeKeyboard(meta)

		id, err := p.tg.SendMessage(meta.ChatID, response.Text, response.Options)
		if err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}

		return p.answers.SaveMessageID(id, meta)
	}

	return nil
}

func (p *Processor) proceedAnswerError(err error, meta entity.Meta) error {
	if errors.Is(err, answers.ErrWrongOption) {
		_, _ = p.tg.SendMessage(meta.ChatID, fmt.Sprintf("Некорректный ответ: %s", err.Error()), nil)
		return nil
	} else if errors.Is(err, answers.ErrFinishSession) {
		p.removeKeyboard(meta)

		props, err := p.answers.FinishSession(meta)
		if err != nil {
			return fmt.Errorf("can't finish session: %w", err)
		}

		msg, err := p.commands.NewOrder(props, meta)
		if err != nil {
			return fmt.Errorf("can't create order: %w", err)
		}

		_, _ = p.tg.SendMessage(meta.ChatID, msg.Text, msg.Options)
		return nil
	} else if errors.Is(err, answers.ErrNoSession) {
		_, _ = p.tg.SendMessage(meta.ChatID, err.Error(), nil)
		return nil
	}
	return fmt.Errorf("can't handle answer: %w", err)
}

func (p *Processor) removeKeyboard(meta entity.Meta) {
	if meta.CallbackID != "" {
		return
	}
	id, err := p.answers.GetMessageID(meta)
	if err == nil && id != 0 {
		if err := p.tg.RemoveKeyboard(meta.ChatID, id); err != nil {
			log.Printf("can't delete message: %v", err)
		}
	}
}
