package entity

import (
	"fmt"
	"strings"
)

type Order struct {
	ID         int
	UserID     int
	Username   string
	Properties map[*State]string
	MinCost    int
	MaxCost    int
}

func OrderToString(order *Order, withUsername bool) string {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("Заказ № %d\n", order.ID))
	if withUsername {
		res.WriteString(fmt.Sprintf("От пользователя @%s\n", order.Username))
	}
	res.WriteString(PropertiesToString(order))
	res.WriteString("\n")
	res.WriteString(fmt.Sprintf("Приблизительная стоимость: %d-%d\n", order.MinCost, order.MaxCost))
	return res.String()
}

func PropertiesToString(order *Order) string {
	var res strings.Builder
	WriteProperty(StateLayoutType, "Тип макета", order, &res)
	WriteProperty(StatePurpose, "Назначение", order, &res)
	WriteProperty(StateScale, "Масштаб", order, &res)
	WriteProperty(StateSize, "Размер", order, &res)
	WriteProperty(StateMaterial, "Материал", order, &res)
	WriteProperty(StateDetails, "Уровень детализации", order, &res)
	WriteProperty(State3DPrint, "3D печать", order, &res)
	WriteProperty(StateLandscape, "Окружающая территория", order, &res)
	WriteProperty(StateDrawings, "Чертежи", order, &res)
	WriteProperty(StateDeadline, "Срок", order, &res)
	WriteProperty(StateDelivery, "Доставка", order, &res)
	return res.String()
}

func WriteProperty(state *State, label string, order *Order, res *strings.Builder) {
	fmt.Fprintf(res, "%s: %s\n", label, order.Properties[state])
}
