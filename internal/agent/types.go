package agent

import "time"

// Task 任务结构体
type Task struct {
	ID        string
	Type      string // codegen/testgen
	Status    string // pending/running/completed/failed
	Input     string
	Language  string
	Result    *TaskResult
	CreatedAt time.Time
}

type TaskResult struct {
	Output   string
	Error    error
	Duration time.Duration
}

// WorkAgent接口
type WorkAgent interface {
	ProcessTask(input string) (*TaskResult, error)
	RetryTask(taskID string) error
	RecoverTask(taskID string) error
}

// WatcherAgent接口
type WatcherAgent interface {
	WatchTask(taskID string) error
}
