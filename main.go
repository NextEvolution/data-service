package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"strings"
)

var redisPool *pool.Pool

func main() {
	killCh := make(chan bool, 1)

	//setup NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to NATS: %s", err))
	}
	defer nc.Close()

	//setup redis pool
	redisPool, err = pool.New("tcp", "localhost:6379", 2)
	if err != nil {
		panic(fmt.Sprintf("Error setting up Redis Pool: %s", err))
	}

	// Simple Async Subscriber
	fmt.Println("listening ...")
	//nc.Subscribe("data-service", func(m *nats.Msg) {
	//	fmt.Printf("Received a message: %s\n", string(m.Data))
	//})

	// recieve data
	nc.Subscribe("dataservice.put.*.lastscan", func(m *nats.Msg) {
		fmt.Printf("Received a message. Subject: %s, Message: %s\n", m.Subject, string(m.Data))

		subjectSplit := strings.Split(m.Subject, ".")
		sellerId := subjectSplit[2]

		conn, err := redisPool.Get()
		if err != nil {
			panic(fmt.Sprintf("Error getting Redis connection: %s", err))
		}
		defer redisPool.Put(conn)

		if conn.Cmd("SET",sellerId, string(m.Data)).Err != nil {
			panic(fmt.Sprintf("Failed to perform REDIS operation"))
		}
	})

	// put data into redis

	<- killCh
}
