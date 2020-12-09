package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Server(ctx context.Context, addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		fmt.Println("hello")
		s.Shutdown(ctx)
	}()
	return s.ListenAndServe()
}

func ServerApp(ctx context.Context, stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello QCon")
	})
	return Server(ctx, ":8080", mux, stop)
}

func ServerDebug(ctx context.Context, stop <-chan struct{}) error {
	return Server(ctx, ":8081", http.DefaultServeMux, stop)
}

func main() {
	stop := make(chan struct{})
	// 一个退出，全部注销退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	//启动一个httpserver 8080
	g.Go(func() error {
		if err := ServerApp(ctx, stop); err != nil {
			cancel()
			return err
		}
		return nil
	})
	//启动另外一个httpserver 8081
	g.Go(func() error {
		if err := ServerDebug(ctx, stop); err != nil {
			cancel()
			return err
		}
		return nil
	})

	// 监听系统信号
	g.Go(func() error {
		signs := make(chan os.Signal, 1)
		signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-signs:
			log.Println("catch system term signal, quit all server in group")
			cancel()
		case <-ctx.Done():
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}

}
