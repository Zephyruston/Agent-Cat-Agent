package agent

type TaskQueue struct {
	queue chan *Task
}

func NewTaskQueue(size int) *TaskQueue {
	return &TaskQueue{
		queue: make(chan *Task, size),
	}
}

func (q *TaskQueue) Enqueue(task *Task) {
	q.queue <- task
}

func (q *TaskQueue) Dequeue() *Task {
	return <-q.queue
}

func (q *TaskQueue) Close() {
	close(q.queue)
}

func (q *TaskQueue) Chan() <-chan *Task {
	return q.queue
}
