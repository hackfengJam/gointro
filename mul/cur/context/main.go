package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-ctx.Done():
			fmt.Println("ctx.Done")
			return ctx.Err()
		}
	}
}
func ExitControl() {
	/*
		在Go1.7发布时，标准库增加了一个context包，
		用来简化对于处理单个请求的多个Goroutine之间与请求域的数据、超时和退出等操作，官方有博文对此做了专门介绍。
		我们可以用context包来重新实现前面的线程安全退出或超时的控制:
	*/
	// 当并发体超时或main主动停止工作者Goroutine时，每个工作者都可以安全退出。
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}

	time.Sleep(time.Second)
	cancel()

	wg.Wait()

}
func main() {
	ExitControl()
}
