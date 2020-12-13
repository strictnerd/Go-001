package main

import (
	"fmt"
	"sync"
)

func main() {
	//fan-out fan-in
	in := gen(7, 3, 4)
	//使用两个goroutine进行获取数据fan out
	d := sq(in)
	e := sq(in)
	// Consume the merged output from c1 and c2.
	/*for n := range merge(d, e) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}*/
	//直接收处理一个，而忽略后续数据，这回导致goroutine泄露
	done := make(chan struct{}, 2)
	defer close(done)
	f := merge(done, d, e)
	fmt.Println(<-f)

}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	defer wg.Done()
	output := func(c <-chan int) {
		for n := range c {
			select {
			case out <- n:
			case <-done:
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func gen(nums ...int) <-chan int {
	out := make(chan int, 3)
	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}

func sq(nums <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range nums {
			out <- n * n
		}
		//如果还没有被取出，就已经关闭，是否有问题？
		close(out)
	}()
	return out
}
