package notify

import "time"

// Notifier接口
type Notifier interface {
	Send(notification *Notification) error
}

type Notification struct {
	Type      string
	TaskID    string
	Message   string
	Timestamp time.Time
}
