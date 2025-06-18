package notify

import "testing"

func TestConsoleNotifier_Send(t *testing.T) {
	n := NewConsoleNotifier()
	err := n.Send(&Notification{Type: "info", TaskID: "t1", Message: "hello"})
	if err != nil {
		t.Fatal(err)
	}
}
