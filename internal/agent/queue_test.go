package agent

import (
	"testing"
)

func TestTaskQueue(t *testing.T) {
	q := NewTaskQueue(2)
	task1 := &Task{ID: "1"}
	task2 := &Task{ID: "2"}
	q.Enqueue(task1)
	q.Enqueue(task2)

	if got := q.Dequeue(); got.ID != "1" {
		t.Errorf("expected 1, got %s", got.ID)
	}
	if got := q.Dequeue(); got.ID != "2" {
		t.Errorf("expected 2, got %s", got.ID)
	}
}
