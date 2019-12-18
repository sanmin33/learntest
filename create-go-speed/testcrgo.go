package main

import (
	"flag"
	"fmt"
	//	"sync"
	"time"
)

//因为创建协程要比加锁和使用共享变量慢很多，而且协程创建后马上就退出了
//所以count加不加锁对运行效率没影响
func crgo(count *int64, stopCh chan int) {
	//	var mu sync.Mutex
	for {
		select {
		case <-stopCh:
			return
		default:
			go func() {
				//				mu.Lock()
				*count++
				//				mu.Unlock()
			}()
		}
	}
}

func testCrGo(n int) {
	counts := make([]int64, n)
	stopCh := make(chan int)
	for i := 0; i < n; i++ {
		go crgo(&(counts[i]), stopCh)
	}
	time.Sleep(1 * time.Second)
	close(stopCh)
	var total int64 = 0
	for _, i := range counts {
		total += i
	}
	fmt.Println("创建者数量:", n, "创建的协程数量:", total)
}

var crNums = flag.Int("crNums", 4, "创建go协程的创建者数量")

func main() {
	flag.Parse()
	time.Sleep(1 * time.Second)
	for i := 1; i <= *crNums; i++ {
		testCrGo(i)
	}
}
