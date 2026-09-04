package entity

type Command struct {
	Text        string
	Description string
}

var (
	StartCmd = &Command{
		Text:        "/start",
		Description: "Запустить бота",
	}
	HelpCmd = &Command{
		Text:        "/help",
		Description: "Вывести подсказку",
	}
	NewOrderCmd = &Command{
		Text:        "/new_order",
		Description: "Создать новый заказ",
	}
	GetOrdersCmd = &Command{
		Text:        "/get_orders",
		Description: "Посмотреть существующие заказы",
	}
	OrdersCmd = &Command{
		Text: "/orders",
	}
)

var Commands = []*Command{StartCmd, HelpCmd, NewOrderCmd, GetOrdersCmd}
