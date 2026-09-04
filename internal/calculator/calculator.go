package calculator

import (
	"layouts-orders-bot/internal/entity"
	"math"
	"strconv"
	"strings"
)

type Calculator struct{}

func New() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Calculate(props map[entity.State]string) (int, int) {
	minCost := 4500.
	maxCost := 4500.

	switch props[entity.StateLayoutType] {
	case "интерьерный":
		minCost *= 1.15
		maxCost *= 1.25
	}

	switch props[entity.StatePurpose] {
	case "в школу/детский сад":
		minCost *= 0.9
		maxCost *= 0.9
	case "выставочный":
		minCost *= 1.2
		maxCost *= 1.4
	case "подарочный":
		minCost *= 1.05
		maxCost *= 1.05
	}

	k := estimateSize(props[entity.StateSize])
	minCost *= k
	maxCost *= k

	switch props[entity.StateMaterial] {
	case "пенокартон":
		minCost *= 2
		maxCost *= 3
	case "пластик":
		minCost *= 4
		maxCost *= 5
	case "комбинированный":
		minCost *= 2
		maxCost *= 3
	}

	switch props[entity.StateDetails] {
	case "базовая":
		minCost *= 0.8
		maxCost *= 1
	case "высокая":
		minCost *= 1.2
		maxCost *= 1.3
	}

	if props[entity.State3DPrint] == "да" {
		minCost += 1000
		maxCost += 7000
	}

	if props[entity.StateLandscape] == "да" {
		minCost += 500
		maxCost += 3000
	}

	switch props[entity.StateDrawings] {
	case "нет, только идея":
		minCost += 2500
		maxCost += 2500
	case "да, но не полностью":
		minCost += 2500
		maxCost += 2500
	}

	switch props[entity.StateDeadline] {
	case "до 2 дней":
		minCost *= 1.7
		maxCost *= 2
	case "7 дней":
		minCost *= 1.2
		maxCost *= 1.3
	}

	return int(math.Floor(minCost/1000)) * 1000, int(math.Ceil(maxCost/1000)) * 1000
}

func estimateSize(size string) float64 {
	num1 := strings.Split(size, "x")[0]
	num2 := strings.Split(size, "x")[1]

	N, _ := strconv.Atoi(num1)
	M, _ := strconv.Atoi(num2)

	return math.Sqrt(float64(N*M) / (15 * 15))
}
