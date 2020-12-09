1.基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。

代码大概是这样实现的：
1. 使用根context创建一个协程
2. 利用1中创建的根context创建子协程（两个http服务和一个信号监听协程）
3. 出现关闭信号调用cancel，Done接收到关闭信号，关闭http服务。

这样就达到了关闭服务，信号协程收到信号，关闭http服务，输出如下日志：
```
stop
stop
2020/12/09 15:03:28 catch system term signal, quit all server in group
2020/12/09 15:03:28 http: Server closed
```


ref:https://golang.org/doc/code.html

https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html
> 本篇精华:如果你不知道一个goroutine何时停止，那么就不要创建这个goroutine.本篇刚开始直接创建了一个永远不会停止的goroutine，很明显这种goroutine会导致泄露，
>然后创建了一个不带缓冲区的channel，并且为这个channel创建了一个超时context，最后导致超时无法正常完成数据接收，最后创建了一个带缓冲区的channel，通过这个channel可以实现goroutne之间的数据通信。

https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html
>保证在服务停止之前所有的goroutine已经关闭，主要是使用go中提供的sync.WaitGroup，可以一直等，也可以设置超时。

https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html
>通过设置gomaxprocs，并结合sync.WaitGroup控制并发和并行

https://dave.cheney.net/practical-go/presentations/qcon-china.html#_concurrency
>最优雅的使用goroutine是交给调用者自己控制，不要过度使用goroutine，最后作者通过举例说明如何控制两个httpserver服务的优雅开启和关闭作为结束。