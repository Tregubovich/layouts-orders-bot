package entity

type Type int

const (
	EventTypeUnknown Type = iota
	EventTypeCommand
	EventTypeAnswer
)

type Event struct {
	Type Type
	Data string
	Meta Meta
}

type Meta struct {
	CallbackID string
	ChatID     int
	Username   string
	UserID     int
}
