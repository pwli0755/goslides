package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	baseCtx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(baseCtx)

	// 启动第一个 HTTP 服务
	g.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello from server 1")
		})
		server := &http.Server{Addr: ":8080", Handler: mux}
		go func() {
			<-ctx.Done()
			server.Shutdown(context.Background())
		}()
		return server.ListenAndServe()
	})

	// 启动第二个 HTTP 服务
	g.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello from server 2")
		})
		server := &http.Server{Addr: ":8081", Handler: mux}
		go func() {
			<-ctx.Done()
			server.Shutdown(context.Background())
		}()
		return server.ListenAndServe()
	})

	// 监听系统信号，优雅地关闭服务
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	g.Go(func() error {
		select {
		case <-c:
			// 收到退出信号，取消所有任务
			fmt.Println("收到退出信号，开始优雅地关闭服务")
			// 通知所有任务退出
			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})

	// 等待所有任务退出
	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("HTTP 服务启动出现错误:", err)
	} else {
		fmt.Println("所有 HTTP 服务优雅地关闭")
	}
}
