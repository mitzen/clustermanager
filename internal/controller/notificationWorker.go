package controller

type NotificationWorker struct {
	messenger Messenger
}

func NewNotificationWorker(messenger Messenger) *NotificationWorker {
	return &NotificationWorker{messenger: messenger}
}

func (s *NotificationWorker) SendMessage(message string) int16 {
	return s.messenger.SendMessage(message)
}
