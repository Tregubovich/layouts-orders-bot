package answers

import (
	"layouts-orders-bot/internal/entity"
	"strconv"
	"strings"
)

const (
	AcceptMessage = "подтвердить"
	CancelMessage = "отмена"
)

type Question struct {
	State           entity.State
	Text            string
	Options         []string
	SpecialValidate func(string) (string, error)
}

var Questions = []*Question{
	{
		State:           entity.StateLayoutType,
		Text:            "Выберите тип макета",
		Options:         []string{"архитектурный", "интерьерный", "промышленный", "ландшафтный"},
		SpecialValidate: nil,
	},
	{
		State:           entity.StatePurpose,
		Text:            "Какое назначение у макета?",
		Options:         []string{"учебный/студенческий", "в школу/детский сад", "выставочный", "подарочный"},
		SpecialValidate: nil,
	},
	{
		State:   entity.StateScale,
		Text:    "Выберите масштаб макета или введите свой",
		Options: []string{"1:50", "1:100", "1:200"},
		SpecialValidate: func(answer string) (string, error) {
			errMsg := &WrongOptionError{"масштаб должен быть в формате '1:X'"}

			data := strings.Split(answer, ":")
			if len(data) != 2 || data[0] != "1" {
				return "", errMsg
			}
			_, err := strconv.Atoi(data[1])
			if err != nil {
				return "", errMsg
			}
			return answer, nil
		},
	},
	{
		State:   entity.StateSize,
		Text:    "Выберите ориентировочный размер макета в см или введите свой",
		Options: []string{"15x15", "30x30", "50x50", "100x100"},
		SpecialValidate: func(answer string) (string, error) {
			errMsg := &WrongOptionError{"размер должен быть в формате 'NxN'"}

			data := strings.Split(answer, "x")
			if len(data) != 2 {
				return "", errMsg
			}
			_, err := strconv.Atoi(data[0])
			if err != nil {
				return "", errMsg
			}
			_, err = strconv.Atoi(data[1])
			if err != nil {
				return "", errMsg
			}
			return answer, nil
		},
	},
	{
		State:           entity.StateMaterial,
		Text:            "Выберите основной материал макета",
		Options:         []string{"пластик", "картон", "другой"},
		SpecialValidate: nil,
	},
	{
		State:           entity.StateDetails,
		Text:            "Выберите детализацию макета",
		Options:         []string{"базовая", "средняя", "высокая"},
		SpecialValidate: nil,
	},
	{
		State:           entity.State3DPrint,
		Text:            "Нужна ли 3D печать?",
		Options:         []string{"да", "нет"},
		SpecialValidate: nil,
	},
	{
		State:           entity.StateLandscape,
		Text:            "Нужна ли окружающая территория?",
		Options:         []string{"да", "нет"},
		SpecialValidate: nil,
	},
	{
		State:           entity.StateDrawings,
		Text:            "Есть ли чертежи?",
		Options:         []string{"нет, только идея", "да, но не полностью", "да, с подписанными размерами"},
		SpecialValidate: nil,
	},
	{
		State:           entity.StateDeadline,
		Text:            "Срок через сколько нужен макет",
		Options:         []string{"до 2 дней", "7 дней", "14 дней", "месяц", "> месяца"},
		SpecialValidate: nil,
	},
	{
		State:           entity.StateDelivery,
		Text:            "В каком городе вы находитесь?",
		Options:         []string{"москва", "санкт-петербург", "другой"},
		SpecialValidate: nil,
	},
	{
		State:   entity.StateCost,
		Text:    "Введите минимальную и максимальную стоимость",
		Options: []string{},
		SpecialValidate: func(answer string) (string, error) {
			errMsg := &WrongOptionError{"стоимость должна быть в формате X-Y"}

			data := strings.Split(answer, "-")
			num1 := strings.TrimSpace(data[0])
			num2 := strings.TrimSpace(data[1])
			if len(data) != 2 {
				return "", errMsg
			}
			_, err := strconv.Atoi(num1)
			if err != nil {
				return "", errMsg
			}
			_, err = strconv.Atoi(num2)
			if err != nil {
				return "", errMsg
			}
			return num1 + "-" + num2, nil
		},
	},
	{
		State:           entity.StateAccept,
		Text:            "%s\nПриблизительная стоимость: %d-%d\n\nПодтвердить заказ?",
		Options:         []string{AcceptMessage, CancelMessage},
		SpecialValidate: nil,
	},
}

func FromStringToOptions(options []string) [][]entity.Option {
	res := make([][]entity.Option, len(options))
	for i, option := range options {
		res[i] = []entity.Option{{Text: option, Data: option}}
	}
	return res
}
