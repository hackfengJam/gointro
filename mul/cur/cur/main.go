package cur

import (
	"fmt"
	"sync"
	"time"
)

func TimeOut() {
	/*
		基于select实现的管道的超时判断：
	*/
	in := make(chan int)
	select {
	case v := <-in:
		fmt.Println(v)
	case <-time.After(time.Second):
		// 超时
		fmt.Println("Time Out !!!")
		return
	}
}

func NoData() {
	/*
		通过select的default分支实现非阻塞的管道发送或接收操作：
	*/
	in := make(chan int)
	select {
	case v := <-in:
		fmt.Println(v)
	default:
		// 没有数据
		fmt.Println("no data !!!")
		return
	}
}

func RandomSeq() {
	/*
		当有多个管道均可操作时，select会随机选择一个管道。基于该特性我们可以用select实现一个生成随机数序列的程序：


	*/
	ch := make(chan int)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
		}
	}()
	for v := range ch {
		fmt.Println(v)
	}
}

func worker(cancel chan bool) {
	// 我们通过select和default分支可以很容易实现一个Goroutine的退出控制:
	for {
		select {
		default:
			fmt.Println("hello")
		// 正常工作
		case <-cancel:
			// 退出
		}
	}
}
func ExitControl() {
	// 我们通过select和default分支可以很容易实现一个Goroutine的退出控制:
	cancel := make(chan bool)
	go worker(cancel)

	time.Sleep(time.Second)
	cancel <- true

}

func ExitControl2() {
	/*
		但是管道的发送操作和接收操作是一一对应的，如果要停止多个Goroutine那么可能需要创建同样数量的管道，这个代价太大了。
		其实我们可以通过close关闭一个管道来实现广播的效果，所有从关闭管道接收的操作均会收到一个零值和一个可选的失败标志。

		我们通过close来关闭cancel管道向多个Goroutine广播退出的指令。不过这个程序依然不够稳健：
		当每个Goroutine收到退出指令退出时一般会进行一定的清理工作，但是退出的清理工作并不能保证被完成，
		因为main线程并没有等待各个工作Goroutine退出工作完成的机制。我们可以结合sync.WaitGroup来改进:
	*/
	cancel := make(chan bool)
	go worker(cancel)

	time.Sleep(time.Second)
	close(cancel)

}

func worker2(wg *sync.WaitGroup, cancel chan bool) {
	/*
		我们通过close来关闭cancel管道向多个Goroutine广播退出的指令。不过这个程序依然不够稳健：
		当每个Goroutine收到退出指令退出时一般会进行一定的清理工作，但是退出的清理工作并不能保证被完成，
		因为main线程并没有等待各个工作Goroutine退出工作完成的机制。我们可以结合sync.WaitGroup来改进:
	*/
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
		case <-cancel:
			return
		}
	}

}

func ExitControl3() {
	/*
		我们通过close来关闭cancel管道向多个Goroutine广播退出的指令。不过这个程序依然不够稳健：
		当每个Goroutine收到退出指令退出时一般会进行一定的清理工作，但是退出的清理工作并不能保证被完成，
		因为main线程并没有等待各个工作Goroutine退出工作完成的机制。我们可以结合sync.WaitGroup来改进:
	*/
	// 现在每个工作者并发体的创建、运行、暂停和退出都是在main函数的安全控制之下了。
	cancel := make(chan bool)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker2(&wg, cancel)
	}
	time.Sleep(time.Second)
	close(cancel)
	wg.Wait()
}
func main() {
	//TimeOut()
	//NoData()
	//RandomSeq()
	//ExitControl()
	//ExitControl2()
	ExitControl3()
}
