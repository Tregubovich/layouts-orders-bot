package commands

import "errors"

var (
	ErrNewSession   = errors.New("should start new session")
	ErrGetOrders    = errors.New("should get orders")
	ErrGetAllOrders = errors.New("should get all orders")
)
