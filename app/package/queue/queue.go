package queue

import (
	"fmt"
	"sync"
)

type Queue[T any] struct {
	mu    *sync.Mutex
	queue []T
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		mu:    &sync.Mutex{},
		queue: []T{},
	}
}

func (q *Queue[T]) Push(r T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, r)
	fmt.Printf("%v\n", q.queue)
}

func (q *Queue[T]) Pop() T {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		var result T
		return result
	}
	r := q.queue[0]
	q.queue = q.queue[1:]
	return r
}
