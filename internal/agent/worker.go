package agent

import (
	"context"
	"time"
)

type Worker struct {
	queue *TaskQueue
	quit  chan struct{}
}

func NewWorker(q *TaskQueue) *Worker {
	return &Worker{queue: q, quit: make(chan struct{})}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case task := <-w.queue.Chan():
				task.SetStatus(StatusRunning)
				// 模拟任务处理
				time.Sleep(10 * time.Millisecond)
				task.SetStatus(StatusCompleted)
			case <-w.quit:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.quit)
}
