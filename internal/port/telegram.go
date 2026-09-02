package port

import (
	"layouts-orders-bot/internal/clients/client"
	"layouts-orders-bot/internal/entity"
	"strings"
)

func EventFromUpdate(upd client.Update) entity.Event {
	if upd.Message != nil {
		upd.Message.Text = strings.ToLower(strings.TrimSpace(upd.Message.Text))

		var eventType entity.Type
		if strings.HasPrefix(upd.Message.Text, "/") {
			eventType = entity.EventTypeCommand
		} else {
			eventType = entity.EventTypeAnswer
		}
		return entity.Event{
			Type: eventType,
			Data: upd.Message.Text,
			Meta: entity.Meta{
				ChatID:   upd.Message.Chat.ID,
				Username: upd.Message.From.Username,
				UserID:   upd.Message.From.ID,
			},
		}
	} else if upd.Callback != nil {
		upd.Callback.Data = strings.ToLower(strings.TrimSpace(upd.Callback.Data))

		return entity.Event{
			Type: entity.EventTypeAnswer,
			Data: upd.Callback.Data,
			Meta: entity.Meta{
				CallbackID: upd.Callback.ID,
				ChatID:     upd.Callback.Message.Chat.ID,
				Username:   upd.Callback.From.Username,
				UserID:     upd.Callback.From.ID,
			},
		}
	} else {
		return entity.Event{
			Type: entity.EventTypeUnknown,
		}
	}
}
