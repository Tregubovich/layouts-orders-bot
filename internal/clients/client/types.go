package client

import "layouts-orders-bot/internal/entity"

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID       int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
	Callback *CallbackQuery   `json:"callback_query"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	Data    string           `json:"data"`
	From    From             `json:"from"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func FromArrayToMarkup(keyboard [][]entity.Option) InlineKeyboardMarkup {
	res := InlineKeyboardMarkup{
		Keyboard: make([][]InlineKeyboardButton, 0, len(keyboard)),
	}
	for i := range keyboard {
		res.Keyboard = append(res.Keyboard, make([]InlineKeyboardButton, 0, len(keyboard[i])))
		for j := range keyboard[i] {
			res.Keyboard[i] = append(res.Keyboard[i], InlineKeyboardButton{
				Text:         keyboard[i][j].Text,
				CallbackData: keyboard[i][j].Data,
			})
		}
	}
	return res
}
