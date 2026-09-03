package entity

type Option struct {
	Text string
	Data string
}

type Message struct {
	Text    string
	Options [][]Option
}
