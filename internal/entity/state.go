package entity

import (
	"errors"
	"strconv"
	"strings"
)

type State struct {
	ID                string
	Options           []string
	SpecialValidation func(string) (string, error)
}

func (s *State) GetOptions() [][]Option {
	res := make([][]Option, len(s.Options))
	for i, option := range s.Options {
		res[i] = []Option{{Text: option, Data: option}}
	}
	return res
}

const (
	LayoutTypeArchitectural = "архитектурный"
	LayoutTypeInterior      = "интерьерный"
	LayoutTypeIndustrial    = "промышленный"
	LayoutTypeLandscape     = "ландшафтный"
)

var StateLayoutType = &State{
	ID: "layout_type",
	Options: []string{
		LayoutTypeArchitectural,
		LayoutTypeInterior,
		LayoutTypeIndustrial,
		LayoutTypeLandscape,
	},
}

const (
	PurposeEducational = "учебный/студенческий"
	PurposeSchool      = "в школу/детский сад"
	PurposeExhibition  = "выставочный"
	PurposeGift        = "подарочный"
)

var StatePurpose = &State{
	ID: "purpose",
	Options: []string{
		PurposeEducational,
		PurposeSchool,
		PurposeExhibition,
		PurposeGift,
	},
}

const (
	Scale1To50  = "1:50"
	Scale1To100 = "1:100"
	Scale1To200 = "1:200"
)

var StateScale = &State{
	ID: "scale",
	Options: []string{
		Scale1To50,
		Scale1To100,
		Scale1To200,
	},
	SpecialValidation: func(answer string) (string, error) {
		errMsg := errors.New("масштаб должен быть в формате '1:X'")

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
}

const (
	Size15x15   = "15x15"
	Size30x30   = "30x30"
	Size50x50   = "50x50"
	Size100x100 = "100x100"
)

var StateSize = &State{
	ID: "size",
	Options: []string{
		Size15x15,
		Size30x30,
		Size50x50,
		Size100x100,
	},
	SpecialValidation: func(answer string) (string, error) {
		errMsg := errors.New("размер должен быть в формате 'NxM'")

		var data []string
		splitSyms := []string{"x", "х", "*", "×", "на"}
		for _, sym := range splitSyms {
			if strings.Contains(answer, sym) {
				data = strings.Split(answer, sym)
				break
			}
		}
		if len(data) != 2 {
			return "", errMsg
		}
		num1 := strings.TrimSpace(data[0])
		num2 := strings.TrimSpace(data[1])
		_, err := strconv.Atoi(num1)
		if err != nil {
			return "", errMsg
		}
		_, err = strconv.Atoi(num2)
		if err != nil {
			return "", errMsg
		}
		return num1 + "x" + num2, nil
	},
}

const (
	MaterialCardboard = "картон"
	MaterialFoamBoard = "пенокартон"
	MaterialPlastic   = "пластик"
	MaterialCombined  = "комбинированный"
)

var StateMaterial = &State{
	ID: "material",
	Options: []string{
		MaterialCardboard,
		MaterialFoamBoard,
		MaterialPlastic,
		MaterialCombined,
	},
}

const (
	DetailsBasic  = "базовая"
	DetailsMedium = "средняя"
	DetailsHigh   = "высокая"
)

var StateDetails = &State{
	ID: "details",
	Options: []string{
		DetailsBasic,
		DetailsMedium,
		DetailsHigh,
	},
}

const (
	Print3DYes = "да"
	Print3DNo  = "нет"
)

var State3DPrint = &State{
	ID: "3d_print",
	Options: []string{
		Print3DYes,
		Print3DNo,
	},
}

const (
	LandscapeYes = "да"
	LandscapeNo  = "нет"
)

var StateLandscape = &State{
	ID: "landscape",
	Options: []string{
		LandscapeYes,
		LandscapeNo,
	},
}

const (
	DrawingsNone     = "нет, только идея"
	DrawingsPartial  = "да, но не полностью"
	DrawingsComplete = "да, с подписанными размерами"
)

var StateDrawings = &State{
	ID: "drawings",
	Options: []string{
		DrawingsNone,
		DrawingsPartial,
		DrawingsComplete,
	},
}

const (
	Deadline2Days  = "до 2 дней"
	Deadline7Days  = "7 дней"
	Deadline14Days = "14 дней"
	DeadlineMonth  = "месяц"
	DeadlineLonger = "> месяца"
)

var StateDeadline = &State{
	ID: "deadline",
	Options: []string{
		Deadline2Days,
		Deadline7Days,
		Deadline14Days,
		DeadlineMonth,
		DeadlineLonger,
	},
}

const (
	DeliveryMoscow = "москва"
	DeliverySPB    = "санкт-петербург"
	DeliveryOther  = "другой"
)

var StateDelivery = &State{
	ID: "delivery",
	Options: []string{
		DeliveryMoscow,
		DeliverySPB,
		DeliveryOther,
	},
}

const (
	AcceptMessage = "Подтвердить"
	CancelMessage = "Отменить"
)

var StateAccept = &State{
	ID: "accept",
	Options: []string{
		AcceptMessage,
		CancelMessage,
	},
}
