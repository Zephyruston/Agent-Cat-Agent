package agent

type DefaultWatcherAgent struct{}

func NewWatcherAgent() *DefaultWatcherAgent {
	return &DefaultWatcherAgent{}
}

func (w *DefaultWatcherAgent) WatchTask(taskID string) error {
	// TODO: 实现任务状态轮询与健康检查
	return nil
}
