package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, c := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return serverApp()
	})

	g.Go(func() error {
		return serverPProf()
	})

	g.Go(func() error {
		HandleSignal(c)
		return nil
	})

	err := g.Wait()
	if err != nil {
		fmt.Println("进程终止", err)
	}
}

func serverApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "hello world")
	})
	return http.ListenAndServe("127.0.0.1:8080", mux)
}
func serverPProf() error {
	return http.ListenAndServe("127.0.0.1:8081", http.DefaultServeMux)
}

func HandleSignal(ctx context.Context) {
	csg := make(chan os.Signal)
	signal.Notify(csg, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range csg {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("正在退出:", s)
			select {
			case <-ctx.Done():
				os.Exit(0)
			case <-Clean():
				fmt.Println("清理完毕...")
				os.Exit(0)
			//设置超时退出
			case <-time.After(time.Second * 1):
				fmt.Println("清理超时...")
				os.Exit(0)
			}
		default:
			fmt.Println("其他信号:", s)
		}
	}
}

func Clean() chan int {
	c := make(chan int)
	fmt.Println("执行清理中...")
	//模拟超时情况
	go func() {
		t := rand.Int31n(3)
		fmt.Println(t, "秒后清理完毕")
		time.Sleep(time.Second * time.Duration(t))
		c <- 1
	}()
	return c
}
