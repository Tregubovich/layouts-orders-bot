package answers

import (
	"layouts-orders-bot/internal/entity"
)

type Question struct {
	State *entity.State
	Text  string
}

var Questions = []*Question{
	{
		State: entity.StateLayoutType,
		Text:  "Выберите тип макета",
	},
	{
		State: entity.StatePurpose,
		Text:  "Какое назначение у макета?",
	},
	{
		State: entity.StateScale,
		Text:  "Выберите масштаб макета или введите свой",
	},
	{
		State: entity.StateSize,
		Text:  "Выберите ориентировочный размер макета в см или введите свой",
	},
	{
		State: entity.StateMaterial,
		Text:  "Выберите основной материал макета",
	},
	{
		State: entity.StateDetails,
		Text:  "Выберите детализацию макета",
	},
	{
		State: entity.State3DPrint,
		Text:  "Нужна ли 3D печать?",
	},
	{
		State: entity.StateLandscape,
		Text:  "Нужна ли окружающая территория?",
	},
	{
		State: entity.StateDrawings,
		Text:  "Есть ли чертежи?",
	},
	{
		State: entity.StateDeadline,
		Text:  "Срок через сколько нужен макет",
	},
	{
		State: entity.StateDelivery,
		Text:  "В каком городе вы находитесь?",
	},
	{
		State: entity.StateAccept,
		Text:  "%s\nПриблизительная стоимость: %d-%d\n\nПодтвердить заказ?",
	},
}
