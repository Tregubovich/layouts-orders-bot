package commands

import (
	"errors"
	"layouts-orders-bot/internal/entity"
	"layouts-orders-bot/internal/handlers/commands/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	meta = entity.Meta{
		ChatID:   123,
		UserID:   456,
		Username: "some_user",
	}

	sampleProps = map[entity.State]string{
		"layout_type": "архитектурный",
		"purpose":     "учебный/студенческий",
		"scale":       "1:100",
		"size":        "30x30",
		"material":    "пластик",
		"details":     "базовая",
		"3d_print":    "да",
		"landscape":   "нет",
		"drawings":    "да, с подписанными размерами",
		"deadline":    "7 дней",
		"delivery":    "москва",
	}

	orders = []*entity.Order{
		{
			ID:         12948192,
			UserID:     meta.UserID,
			Properties: sampleProps,
			MinCost:    12000,
			MaxCost:    25000,
		},
		{
			ID:         12948283,
			UserID:     meta.UserID,
			Properties: sampleProps,
			MinCost:    11000,
			MaxCost:    20000,
		},
		{
			ID:         12948283,
			UserID:     meta.UserID + 12,
			Properties: sampleProps,
			MinCost:    1000,
			MaxCost:    20500,
		},
	}
)

func TestHandleCommand(t *testing.T) {
	handler, _, _ := setup(t)

	msg, err := handler.HandleCommand(entity.StartCmd.Text, meta)
	require.NoError(t, err)
	require.Equal(t, msgStart, msg.Text)

	msg, err = handler.HandleCommand(entity.HelpCmd.Text, meta)
	require.NoError(t, err)
	for _, c := range entity.Commands {
		require.Contains(t, msg.Text, c.Text)
		require.Contains(t, msg.Text, c.Description)
	}

	msg, err = handler.HandleCommand(entity.NewOrderCmd.Text, meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNewSession)

	msg, err = handler.HandleCommand(entity.GetOrdersCmd.Text, meta)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrGetOrders)

	msg, err = handler.HandleCommand(entity.OrdersCmd.Text, meta)
	require.NoError(t, err)
	require.Equal(t, msgUnknownCommand, msg.Text)

	msg, err = handler.HandleCommand("/unknown_command", meta)
	require.NoError(t, err)
	require.Equal(t, msgUnknownCommand, msg.Text)
}

func TestNewOrder(t *testing.T) {
	handler, repo, calc := setup(t)

	minCost, maxCost := 12000, 25000
	calc.EXPECT().Calculate(sampleProps).Return(minCost, maxCost)
	repo.EXPECT().NewOrder(gomock.Any()).DoAndReturn(
		func(order *entity.Order) error {
			require.Equal(t, meta.UserID, order.UserID)
			require.Equal(t, sampleProps, order.Properties)
			require.Equal(t, minCost, order.MinCost)
			require.Equal(t, maxCost, order.MaxCost)
			return nil
		})

	msg, err := handler.NewOrder(sampleProps, meta)
	require.NoError(t, err)
	require.Equal(t, msgAccepted, msg.Text)
	require.Nil(t, msg.Options)
}

func TestNewOrderWithError(t *testing.T) {
	handler, repo, calc := setup(t)

	minCost, maxCost := 12000, 25000
	calc.EXPECT().Calculate(sampleProps).Return(minCost, maxCost)
	repo.EXPECT().NewOrder(gomock.Any()).Return(errors.New("some error"))

	msg, err := handler.NewOrder(sampleProps, meta)
	require.Error(t, err)
	require.ErrorContains(t, err, "some error")
	require.Nil(t, msg)
}

func TestGetOrders(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetOrders(meta.UserID).Return(orders[:len(orders)-1], nil)

	messages, err := handler.GetOrders(meta, false)
	require.NoError(t, err)
	require.Len(t, messages, 2)
	for i, msg := range messages {
		require.Equal(t, msg.Text, entity.OrderToString(orders[i], false))
		require.Nil(t, msg.Options)
	}
}

func TestGetAllOrders(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetAllOrders().Return(orders, nil)

	messages, err := handler.GetOrders(meta, true)
	require.NoError(t, err)
	require.Len(t, messages, len(orders))
	for i, msg := range messages {
		require.Equal(t, msg.Text, entity.OrderToString(orders[i], true))
		require.Nil(t, msg.Options)
	}
}

func TestGetOrdersWithError(t *testing.T) {
	handler, repo, _ := setup(t)

	repo.EXPECT().GetOrders(meta.UserID).Return(nil, errors.New("some error"))

	msg, err := handler.GetOrders(meta, false)
	require.Error(t, err)
	require.ErrorContains(t, err, "some error")
	require.Nil(t, msg)
}

func setup(t *testing.T) (*Handler, *mocks.MockRepository, *mocks.MockCalculator) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := mocks.NewMockRepository(ctrl)
	calc := mocks.NewMockCalculator(ctrl)

	return NewHandler(repo, calc), repo, calc
}
