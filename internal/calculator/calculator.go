package calculator

import (
	"layouts-orders-bot/internal/entity"
)

type Calculator struct{}

func New() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Calculate(props map[entity.State]string) (int, int) {
	return 0, 100
}
