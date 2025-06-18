package agent

import "testing"

func TestTaskStatus(t *testing.T) {
	task := &Task{ID: "1"}
	task.SetStatus(StatusPending)
	if task.Status != string(StatusPending) {
		t.Errorf("expected pending, got %s", task.Status)
	}
	task.SetStatus(StatusRunning)
	if task.Status != string(StatusRunning) {
		t.Errorf("expected running, got %s", task.Status)
	}
	task.SetStatus(StatusCompleted)
	if !task.IsTerminal() {
		t.Error("completed should be terminal")
	}
}
