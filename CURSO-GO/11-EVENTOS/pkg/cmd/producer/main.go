package main

import "github.com/dlcdev1/pos1/eventos/utils/pkg/rabbitmq"

func main() {
	ch, err := rabbitmq.OpenChanel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbitmq.Publish(ch, "Hello World!", "amq.direct")
}
