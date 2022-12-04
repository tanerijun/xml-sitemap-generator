package queue

type Queue struct {
	queue []string
}

func (q *Queue) Enqueue(s string) {
	q.queue = append(q.queue, s)
}

func (q *Queue) Top() string {
	return q.queue[0]
}

func (q *Queue) Dequeue() string {
	ret := q.queue[0]
	q.queue = q.queue[1:]
	return ret
}

func (q *Queue) Empty() bool {
	return len(q.queue) == 0
}

func New(items []string) *Queue {
	return &Queue{
		queue: items,
	}
}
