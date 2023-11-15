package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// START
func main() {
	var value int32 = 10

	// 假设有两个 goroutine 同时想要将 value 的值从 10 修改为 20
	go func() {
		// 在这个 goroutine 中，我们尝试将 value 的值从 10 修改为 20
		// 只有当 value 的值仍然是 10 时，才会执行更新操作
		swapped := atomic.CompareAndSwapInt32(&value, 10, 20)
		fmt.Println("Goroutine 1 attempted to swap:", swapped)
	}()

	go func() {
		// 在另一个 goroutine 中，我们也尝试将 value 的值从 10 修改为 20
		// 由于另一个 goroutine 已经执行了更新操作，这里的 CAS 操作将不会成功
		swapped := atomic.CompareAndSwapInt32(&value, 10, 20)
		fmt.Println("Goroutine 2 attempted to swap:", swapped)
	}()

	// 等待两个 goroutine 执行完毕
	// 由于 CAS 操作是原子的，因此只有一个 goroutine 的更新操作会成功
	time.Sleep(1 * time.Second)
}

// END
