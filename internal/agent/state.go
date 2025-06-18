package agent

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

func (t *Task) SetStatus(status TaskStatus) {
	t.Status = string(status)
}

func (t *Task) IsTerminal() bool {
	return t.Status == string(StatusCompleted) || t.Status == string(StatusFailed)
}
