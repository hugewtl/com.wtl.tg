package main

import (
	"fmt"
	"sync"
)

type worker chan struct {
	intChan chan int
	wg      *sync.WaitGroup
}

func main() {
	// //管道要做初始化
	// // var wg1 *sync.WaitGroup
	// workers := worker{
	// 	intChan: make(chan int)
	// 	wg1:     wg

	// }
	// 	intChan = make(chan int, 20)
	// 	// wg.Add(1)
	// 	// go fmt.Println(gorIn(intChan, wg))
	// 	// // time.Sleep(time.Second)
	// 	// // go gorOut(intChan, wg)
	// 	// wg.Wait()
	// }
}
func gorIn(ch chan int, wg *sync.WaitGroup) int {
	for i := 0; i < 1000000000; i++ {
		ch <- i
		// fmt.Println("ch=", <-ch)
		return <-ch
	}
	wg.Done()
	return <-ch
}

func gorOut(ch <-chan int, wg *sync.WaitGroup) {
	for i := 0; i < 1000000000; i++ {
		k := i % <-ch
		fmt.Println("k=", k)
	}
	wg.Done()
}
