package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// START
func main() {
	c := sync.NewCond(&sync.Mutex{})
	var ready atomic.Int32

	for i := 0; i < 3; i++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)

			// 加锁更改等待条件
			ready.Add(1)
			fmt.Printf("运动员#%d 已准备就绪\n", i)
			// 广播唤醒所有的等待者
			c.Broadcast()
		}(i)
	}

	// 调用 cond.Wait 方法之前一定要加锁
	c.L.Lock()
	// waiter goroutine 被唤醒不等于等待条件被满足，只是有 goroutine 把它唤醒了而已
	// 等待条件有可能已经满足了，也有可能不满足，我们需要进一步检查
	for ready.Load() != 3 {
		c.Wait()
		fmt.Println("裁判员被唤醒一次")
	}
	c.L.Unlock()

	//所有的运动员是否就绪
	fmt.Println("所有运动员都准备就绪。比赛开始，3，2，1, ......")
}

// END
