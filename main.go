package main

import (
	"github.com/nats-io/nats"
	"fmt"
)

func main() {
	killCh := make(chan bool, 1)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to NATS: %s", err))
	}

	// Simple Async Subscriber
	fmt.Println("listening ...")
	nc.Subscribe("data-service", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	<- killCh
}
