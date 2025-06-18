# internal/agent

WorkAgent与WatcherAgent核心逻辑实现。

- 任务队列（channel+goroutine）
- Task状态机
- 任务持久化（BoltDB）
- Worker异步处理 