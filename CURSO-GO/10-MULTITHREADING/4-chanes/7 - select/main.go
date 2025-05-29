package main

import (
	"fmt"
	"time"
)

type Message struct {
	id  int
	Msg string
}

func main() {

	c1 := make(chan Message)
	c2 := make(chan Message)

	// RabitMQ
	go func() {
		time.Sleep(time.Second)
		msg := Message{1, "Hello from RabbitMQ"}
		c1 <- msg
	}()

	go func() {
		time.Sleep(time.Second)
		msg := Message{1, "Hello from Kafka"}
		c2 <- msg
	}()

	for i := 0; i < 3; i++ {
		select {
		case msg := <-c1:

			fmt.Printf("Received from RabbitMq: %s\n", msg.Msg)
		case msg := <-c2:
			fmt.Printf("Received from kafka %s\n", msg.Msg)
		case <-time.After(time.Second * 3):
			println("timeout")
		default:
			println("default")
		}
	}

}
