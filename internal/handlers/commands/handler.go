package commands

import (
	"fmt"
	"layouts-orders-bot/internal/entity"
	"log"
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

func (h *Handler) HandleCommand(command string, chatID int, username string) (*entity.Message, error) {
	log.Printf("got new command '%s' from '%s", command, username)

	switch command {
	case StartCmd:
		return &entity.Message{Text: msgStart}, nil
	case HelpCmd:
		return &entity.Message{Text: msgHelp}, nil
	case NewOrderCmd:
		return nil, ErrNewSession
	case GetOrdersCmd:
		return nil, ErrGetOrders
	default:
		return &entity.Message{Text: msgUnknownCommand}, nil
	}
}

func (h *Handler) NewOrder(userID int, username string, props map[entity.State]string) (*entity.Order, error) {
	log.Printf("got new order from '%s", username)

	minCost, maxCost := h.calculator.Calculate(props)
	order := &entity.Order{
		UserID:     userID,
		Properties: props,
		MinCost:    minCost,
		MaxCost:    maxCost,
	}

	err := h.repo.NewOrder(order)
	if err != nil {
		return nil, fmt.Errorf("could not create new order: %w", err)
	}
	return order, nil
}

func (h *Handler) GetOrders(userID int, username string) ([]*entity.Message, error) {
	log.Printf("get orders for '%s", username)

	orders, err := h.repo.GetOrders(userID)
	if err != nil {
		return nil, fmt.Errorf("can't get orders: %w", err)
	}

	res := make([]*entity.Message, 0, len(orders))
	for _, order := range orders {
		res = append(res, &entity.Message{Text: entity.OrderToString(order)})
	}

	return res, nil
}
