package memory

import (
	"context"

	"github.com/streampulse/api/internal/queue"
)

type Queue struct {
	ch chan queue.DealEvent
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{ch: make(chan queue.DealEvent, bufferSize)}
}

func (q *Queue) Publish(ctx context.Context, event queue.DealEvent) error {
	select {
	case q.ch <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) Subscribe() <-chan queue.DealEvent {
	return q.ch
}

func (q *Queue) Close() {
	close(q.ch)
}
