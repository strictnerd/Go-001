package main

import (
	"container/ring"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	limitCount  int        = 10 //限频
	limitBucket int        = 4  //滑动窗口个数
	curCount    int32      = 0  //当前数量
	head        *ring.Ring      //环形队列
)

func main() {
	//初始化滑动窗口
	head = ring.New(limitBucket)
	for i := 0; i < limitBucket; i++ {
		head.Value = 0
		head = head.Next()
	}

	go func() {
		timer := time.NewTicker(1 * time.Second)
		//每隔一秒刷新一次滑动窗口数据
		for range timer.C {
			subCount := int32(0 - head.Value.(int))
			newCount := atomic.AddInt32(&curCount, subCount)

			arr := [6]int{}
			for i := 0; i < limitBucket; i++ {
				arr[i] = head.Value.(int)
				head = head.Next()
			}
			fmt.Println("move subCount,newCount,arr", subCount, newCount, arr)
			head.Value = 0
			head = head.Next()
		}
	}()

	for i := 0; i < 15; i++ {
		go func() {
			handle()
		}()
	}
	time.Sleep(50 * time.Second)
}

func handle() {
	n := atomic.AddInt32(&curCount, 1)
	fmt.Println("handler n:", n)
	if n > int32(limitCount) { // 超出限频
		atomic.AddInt32(&curCount, -1) // add 1 by atomic，业务处理完毕，放回令牌
		fmt.Println("too many request, please try again.")
	} else {
		mu := sync.Mutex{}
		mu.Lock()
		pos := head.Prev()
		val := pos.Value.(int)
		val++
		pos.Value = val
		mu.Unlock()
		time.Sleep(1 * time.Second) // 假设我们的应用处理业务用了1s的时间
		fmt.Println("I can change the world!")
	}
}
