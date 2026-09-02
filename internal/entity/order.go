package entity

type Order struct {
	Username   string
	Properties map[State]string
	MinCost    int
	MaxCost    int
}

func OrderToString(order Order) string {
	return ``
}
