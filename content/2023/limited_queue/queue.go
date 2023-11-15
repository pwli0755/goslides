package main

import (
	"fmt"
	"sync"
	"time"
)

type LimitedQueue struct {
	cond     *sync.Cond
	elements []int
	maxSize  int
}

func NewLimitedQueue(maxSize int) *LimitedQueue {
	return &LimitedQueue{
		cond:     sync.NewCond(&sync.Mutex{}),
		elements: make([]int, 0, maxSize),
		maxSize:  maxSize,
	}
}

func (q *LimitedQueue) Enqueue(value int) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.elements) == q.maxSize {
		q.cond.Wait() // 等待队列不满的通知
	}
	q.elements = append(q.elements, value)
	q.cond.Signal() // 通知队列不为空
}

func (q *LimitedQueue) Dequeue() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.elements) == 0 {
		q.cond.Wait() // 等待队列不空的通知
	}
	value := q.elements[0]
	q.elements = q.elements[1:]
	q.cond.Signal() // 通知队列不满
	return value
}

// START
func main() {
	queue := NewLimitedQueue(2)

	var wg sync.WaitGroup
	wg.Add(2)

	// 生产者goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			queue.Enqueue(i)
			fmt.Println("生产者：入队 <<<", i)
		}
	}()

	// 消费者goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			time.Sleep(2 * time.Second)
			value := queue.Dequeue()
			fmt.Println("消费者：出队 >>>", value)
		}
	}()

	wg.Wait()
}

// END
