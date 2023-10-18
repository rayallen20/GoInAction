package queue

import (
	"context"
	"errors"
	"sync"
)

type SliceQueue struct {
	// queue 队列底层
	queue []any
	// mutex 互斥锁
	mutex sync.Mutex
}

func NewSliceQueue(size int) *SliceQueue {
	q := &SliceQueue{
		queue: make([]any, size),
	}
	return q
}

func (q *SliceQueue) Enqueue(v any) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue = append(q.queue, v)
	return nil
}

// Dequeue 也可以让用户判定queue是否为空
// 并发场景下 事前的判定不一定准确
func (q *SliceQueue) Dequeue() (v any, err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// 队列长度为0 则必然没有元素可读
	// 大多数情况下 队列中都是有元素的 所以后边一定还有一次写请求 所以即使这里加一个读锁 意义不大
	if len(q.queue) == 0 {
		return nil, errors.New("队列没有元素")
	}

	// 返回slice中的第0个元素
	v = q.queue[0]
	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		q.queue = make([]any, 0)
	}
	return v, nil
}

func (q *SliceQueue) Traversal(ctx context.Context, f func(any) error) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for _, v := range q.queue {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}
