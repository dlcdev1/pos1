package main

import "fmt"

// Thread 1
func main() {
	ch := make(chan int)
	go publish(ch)
	read(ch)
}

func read(ch chan int) {
	for i := range ch {
		fmt.Printf("Received %d\n", i)
	}
}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}
