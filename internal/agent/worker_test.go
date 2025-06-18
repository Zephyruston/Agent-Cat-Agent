package agent

import (
	"context"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	q := NewTaskQueue(1)
	worker := NewWorker(q)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker.Start(ctx)
	task := &Task{ID: "w1"}
	q.Enqueue(task)
	time.Sleep(20 * time.Millisecond)
	if task.Status != string(StatusCompleted) {
		t.Errorf("expected completed, got %s", task.Status)
	}
	worker.Stop()
}
