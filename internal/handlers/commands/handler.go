package commands

import (
	"fmt"
	"layouts-orders-bot/internal/entity"
	"log"
	"slices"
	"strings"
)

var adminID = []int{1196093685, 869479338}

//go:generate mockgen -source=handler.go -destination=mocks/handler_mocks.go -package=mocks
type (
	Repository interface {
		NewOrder(order *entity.Order) error
		GetOrders(userID int) ([]*entity.Order, error)
		GetAllOrders() ([]*entity.Order, error)
	}

	Calculator interface {
		Calculate(props map[*entity.State]string) (int, int)
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
	log.Printf("(%s) got new command '%s'", meta.Username, command)

	if command == entity.StartCmd.Text {
		return &entity.Message{Text: msgStart}, nil
	} else if command == entity.HelpCmd.Text {
		return h.helpCmd(), nil
	} else if command == entity.NewOrderCmd.Text {
		return nil, ErrNewSession
	} else if command == entity.GetOrdersCmd.Text {
		return nil, ErrGetOrders
	} else if command == entity.OrdersCmd.Text {
		if !slices.Contains(adminID, meta.UserID) {
			return &entity.Message{Text: msgUnknownCommand}, nil
		}
		return nil, ErrGetAllOrders
	} else {
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

func (h *Handler) NewOrder(props map[*entity.State]string, meta entity.Meta) (*entity.Message, error) {
	log.Printf("(%s) got new order", meta.Username)

	minCost, maxCost := h.calculator.Calculate(props)
	order := &entity.Order{
		UserID:     meta.UserID,
		Username:   meta.Username,
		Properties: props,
		MinCost:    minCost,
		MaxCost:    maxCost,
	}

	err := h.repo.NewOrder(order)
	if err != nil {
		return nil, fmt.Errorf("could not create new order: %w", err)
	}
	return &entity.Message{Text: fmt.Sprintf(msgAccepted, order.ID)}, nil
}

func (h *Handler) GetOrders(meta entity.Meta, isAdmin bool) ([]*entity.Message, error) {
	if !isAdmin {
		log.Printf("(%s) getting orders", meta.Username)
	} else {
		log.Printf("(%s) getting orders (admin mode)", meta.Username)
	}

	var orders []*entity.Order
	var err error
	if !isAdmin {
		orders, err = h.repo.GetOrders(meta.UserID)
	} else {
		orders, err = h.repo.GetAllOrders()
	}
	if err != nil {
		return nil, fmt.Errorf("can't get orders: %w", err)
	}

	if len(orders) == 0 {
		return []*entity.Message{{Text: msgNoOrders}}, nil
	}

	res := make([]*entity.Message, 0, len(orders))
	for _, order := range orders {
		res = append(res, &entity.Message{Text: entity.OrderToString(order, isAdmin)})
	}

	return res, nil
}
