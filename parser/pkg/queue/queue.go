package queue

import (
	"sync"
)

type Queue[T any] interface {
	Push(values ...T)
	Pop() (T, bool)
}

type queue[T any] struct {
	ch       chan T
	capacity int
	length   int

	mx sync.RWMutex
}

func NewQueue[T any](capacity int) Queue[T] {
	return &queue[T]{
		ch:       make(chan T, capacity),
		capacity: capacity,
		length:   0,
	}
}

func (q *queue[T]) Push(values ...T) {
	length := q.length + len(values)
	if length > q.capacity {
		q.extend(length)
	}

	for _, value := range values {
		q.ch <- value
	}

	q.length = length
}

func (q *queue[T]) newCap(length int) int {
	capacity := q.capacity
	doubleCapacity := 2 * capacity
	if length > doubleCapacity {
		return length
	}

	const threshold = 256
	if q.capacity < threshold {
		return doubleCapacity
	}

	for uint(capacity) < uint(length) {
		capacity += (capacity + 3*threshold) / 4
	}

	if capacity <= 0 {
		return length
	}

	return capacity
}

func (q *queue[T]) extend(length int) {
	q.mx.Lock()
	defer q.mx.Unlock()

	capacity := q.newCap(length)
	ch := make(chan T, capacity)

	close(q.ch)
	for value := range q.ch {
		ch <- value
	}

	q.ch = ch
	q.capacity = capacity
}

func (q *queue[T]) Pop() (T, bool) {
	var value T
	if q.length == 0 {
		return value, false
	}

	q.mx.RLock()
	defer q.mx.RUnlock()

	var ok bool
	value, ok = <-q.ch
	if ok {
		q.length--
	}

	return value, ok
}
