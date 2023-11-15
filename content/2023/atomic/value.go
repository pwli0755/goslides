package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// START
type ComplexStruct struct {
	counter int
	// 其他复杂的字段
}

func main() {
	var wg sync.WaitGroup
	var complexValue atomic.Value
	complexValue.Store(&ComplexStruct{counter: 0})

	wg.Add(2)
	// 使用 Store 方法在一个 goroutine 中设置共享结构体的值
	go func() {
		defer wg.Done()
		newValue := &ComplexStruct{counter: 100}
		complexValue.Store(newValue)
	}()

	// 使用 Load 方法在另一个 goroutine 中读取共享结构体的值
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		loadedValue := complexValue.Load().(*ComplexStruct)
		fmt.Println("The counter is:", loadedValue.counter)
	}()

	wg.Wait()
}

// END
