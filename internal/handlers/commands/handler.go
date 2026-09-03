package commands

import (
	"fmt"
	"layouts-orders-bot/internal/entity"
	"log"
	"strings"
)

const (
	StartCmd     = "/start"
	HelpCmd      = "/help"
	NewOrderCmd  = "/new_order"
	GetOrdersCmd = "/get_orders"
)

type (
	Repository interface {
		NewOrder(order *entity.Order) error
		GetOrders(userID int) ([]*entity.Order, error)
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

func (h *Handler) HandleCommand(command string, meta entity.Meta) (*entity.Message, error) {
	log.Printf("got new command '%s' from '%s", command, meta.Username)

	switch command {
	case StartCmd:
		return &entity.Message{Text: msgStart}, nil
	case HelpCmd:
		return h.helpCmd(), nil
	case NewOrderCmd:
		return nil, ErrNewSession
	case GetOrdersCmd:
		return nil, ErrGetOrders
	default:
		return &entity.Message{Text: msgUnknownCommand}, nil
	}
}

func (h *Handler) helpCmd() *entity.Message {
	var b strings.Builder
	for _, cmd := range entity.Commands {
		fmt.Fprintf(&b, "%s — %s\n", cmd.Text, cmd.Description)
	}
	return &entity.Message{Text: fmt.Sprintf(msgHelp, b.String())}
}

func (h *Handler) NewOrder(props map[entity.State]string, meta entity.Meta) (*entity.Message, error) {
	log.Printf("got new order from '%s", meta.Username)

	minCost, maxCost := h.calculator.Calculate(props)
	order := &entity.Order{
		UserID:     meta.UserID,
		Properties: props,
		MinCost:    minCost,
		MaxCost:    maxCost,
	}

	err := h.repo.NewOrder(order)
	if err != nil {
		return nil, fmt.Errorf("could not create new order: %w", err)
	}
	return &entity.Message{Text: msgAccepted}, nil
}

func (h *Handler) GetOrders(meta entity.Meta) ([]*entity.Message, error) {
	log.Printf("get orders for '%s", meta.Username)

	orders, err := h.repo.GetOrders(meta.UserID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders: %w", err)
	}

	res := make([]*entity.Message, 0, len(orders))
	for _, order := range orders {
		res = append(res, &entity.Message{Text: entity.OrderToString(order)})
	}

	return res, nil
}
