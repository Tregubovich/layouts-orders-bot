package entity

import (
	"fmt"
	"strings"
)

type Order struct {
	Username   string
	Properties map[State]string
	MinCost    int
	MaxCost    int
}

func OrderToString(order *Order) string {
	props := make([]string, len(order.Properties))

	for state, value := range order.Properties {
		props = append(props, fmt.Sprintf("%v=%v", state, value))
	}

	return strings.Join(props, "\n")
}
