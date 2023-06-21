package controller

type Messenger interface {
	SendMessage(message string) int16
}
