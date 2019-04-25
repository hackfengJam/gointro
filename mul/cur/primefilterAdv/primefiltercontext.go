package main

import (
	"context"
	"fmt"
)

// 返回生成自然数序列的管道: 2, 3, 4, ...
func GenerateNatural(ctx context.Context) chan int {
	ch := make(chan int)

	go func() {
		for i := 2; ; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()

	return ch
}

// 管道过滤器: 删除能被素数整除的数
func PrimeFilter(ctx context.Context, in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; i%prime != 0 {
				select {
				case <-ctx.Done():
					return
				case out <- i:
				}
			}
		}

	}()
	return out
}

func main() {
	// 通过 Context 控制后台 Goroutine状态
	/*
		当main函数完成工作前，通过调用cancel()来通知后台Goroutine退出，这样就避免了Goroutine的泄漏。

		并发是一个非常大的主题，我们这里只是展示几个非常基础的并发编程的例子。
		官方文档也有很多关于并发编程的讨论，国内也有专门讨论Go语言并发编程的书籍。
	*/
	ctx, cancel := context.WithCancel(context.Background())

	ch := GenerateNatural(ctx)
	for i := 0; i < 100; i++ {
		prime := <-ch // 新出现的素数
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilter(ctx, ch, prime)
	}
	cancel()
}
