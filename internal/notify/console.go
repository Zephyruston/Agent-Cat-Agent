package notify

import (
	"fmt"
)

type ConsoleNotifier struct{}

func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

func (n *ConsoleNotifier) Send(notification *Notification) error {
	fmt.Printf("[%s][%s] %s\n", notification.Type, notification.TaskID, notification.Message)
	return nil
}
