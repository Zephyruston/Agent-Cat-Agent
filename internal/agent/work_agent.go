package agent

import (
	"fmt"
)

type DefaultWorkAgent struct {
	queue *TaskQueue
}

func NewWorkAgent(queueSize int) *DefaultWorkAgent {
	return &DefaultWorkAgent{
		queue: NewTaskQueue(queueSize),
	}
}

func (w *DefaultWorkAgent) ProcessTask(input string) (*TaskResult, error) {
	task := &Task{ID: input, Status: string(StatusPending)}
	w.queue.Enqueue(task)
	// 简化：直接返回入队结果，后续可扩展为异步处理
	return &TaskResult{Output: fmt.Sprintf("Task %s enqueued", input)}, nil
}

func (w *DefaultWorkAgent) RetryTask(taskID string) error {
	// TODO: 实现重试逻辑
	return nil
}

func (w *DefaultWorkAgent) RecoverTask(taskID string) error {
	// TODO: 实现恢复逻辑
	return nil
}
