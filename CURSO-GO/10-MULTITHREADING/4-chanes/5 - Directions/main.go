package main

import (
	"fmt"
	"sync"
)

// Thread 1
func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(10)
	go publish(ch)
	go read(ch, &wg)
	wg.Wait()

}

func read(ch chan int, wg *sync.WaitGroup) {
	for i := range ch {
		fmt.Printf("Received %d\n", i)
		wg.Done()
	}
}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}
