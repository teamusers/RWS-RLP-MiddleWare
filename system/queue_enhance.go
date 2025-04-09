package system

import (
	"fmt"
	"sync"

	"rlp-member-service/log"
)

type processHandler[T comparable] func(slot T, lock *sync.WaitGroup)

var signInLocks sync.Map // Thread-safe global lock

type RichQueue[T comparable] struct {
	sync.Mutex
	slots  []T
	notify chan struct{}
}

func NewRichQueue[T comparable]() *RichQueue[T] {
	return &RichQueue[T]{
		slots:  make([]T, 0),
		notify: make(chan struct{}, 1),
	}
}

func (q *RichQueue[T]) Enqueue(item T) {
	_, loaded := signInLocks.LoadOrStore(item, struct{}{})
	if loaded {
		return
	}

	q.Lock()
	defer q.Unlock()

	q.slots = append(q.slots, item)
	select {
	case q.notify <- struct{}{}:
	default:
	}
}

func (q *RichQueue[T]) BatchEnqueue(items []T) {
	q.Lock()
	defer q.Unlock()

	for _, item := range items {
		_, loaded := signInLocks.LoadOrStore(item, struct{}{})
		if loaded {
			continue
		}
		q.slots = append(q.slots, item)
	}
	select {
	case q.notify <- struct{}{}:
	default:
	}
}

func (q *RichQueue[T]) Dequeue() (T, error) {
	q.Lock()
	defer q.Unlock()

	if len(q.slots) == 0 {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}
	item := q.slots[0]
	q.slots = q.slots[1:]

	return item, nil
}

// ✅ BatchDequeue method to prevent returning a nil slice and avoid panic
func (q *RichQueue[T]) BatchDequeue(size int) ([]T, error) {
	q.Lock()
	defer q.Unlock()

	if size <= 0 {
		return []T{}, fmt.Errorf("wrong size") // Return an empty slice
	}
	if len(q.slots) == 0 {
		return []T{}, fmt.Errorf("empty length") // Return an empty slice
	}

	var ret []T
	if len(q.slots) > size {
		ret = make([]T, size)
		copy(ret, q.slots[:size])
		q.slots = q.slots[size:]
	} else {
		ret = q.slots
		q.slots = []T{}
	}
	return ret, nil
}

func (q *RichQueue[T]) Size() int {
	q.Lock()
	defer q.Unlock()
	return len(q.slots)
}

func (q *RichQueue[T]) First() T {
	q.Lock()
	defer q.Unlock()
	var zero T
	if len(q.slots) == 0 {
		return zero
	}
	return q.slots[0]
}

func (q *RichQueue[T]) Last() T {
	q.Lock()
	defer q.Unlock()
	var zero T
	if len(q.slots) == 0 {
		return zero
	}
	return q.slots[len(q.slots)-1]
}

// ✅ Consumer method to avoid deadlock & 100% CPU usage issue
func (q *RichQueue[T]) Consumer(size int, handler processHandler[T]) {
	for {
		<-q.notify
		for {
			items, err := q.BatchDequeue(size)
			if err != nil {
				log.Errorf("Get Consume Batch: %d, %v\n", size, err)
				break // Prevent 100% CPU usage
			}
			var wg sync.WaitGroup
			for _, item := range items {
				wg.Add(1)
				func(item T) {
					defer wg.Done()
					handler(item, &wg)
				}(item)
			}
			wg.Wait()
		}
	}
}
