package entity

import (
	"fmt"
	"strings"
)

type Order struct {
	ID         int
	UserID     int
	Properties map[State]string
	MinCost    int
	MaxCost    int
}

func OrderToString(order *Order) string {
	var res strings.Builder
	AddProperty(StateLayoutType, "Тип макета", order, &res)
	AddProperty(StatePurpose, "Назначение", order, &res)
	AddProperty(StateScale, "Масштаб", order, &res)
	AddProperty(StateSize, "Размер", order, &res)
	AddProperty(StateMaterial, "Материал", order, &res)
	AddProperty(StateDetails, "Уровень детализации", order, &res)
	AddProperty(State3DPrint, "3D печать", order, &res)
	AddProperty(StateLandscape, "Окружающая территория", order, &res)
	AddProperty(StateDrawings, "Чертежи", order, &res)
	AddProperty(StateDeadline, "Срок", order, &res)
	AddProperty(StateDelivery, "Доставка", order, &res)
	return res.String()
}

func AddProperty(state State, label string, order *Order, res *strings.Builder) {
	fmt.Fprintf(res, "%s: %s\n", label, order.Properties[state])
}
