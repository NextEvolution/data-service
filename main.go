package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"strings"
	"github.com/mediocregopher/radix.v2/redis"
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
	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	redisPool, err = pool.NewCustom("tcp", "localhost:6379", 2, df)
	if err != nil {
		panic(fmt.Sprintf("Error setting up Redis Pool: %s", err))
	}

	// Simple Async Subscriber
	fmt.Println("listening ...")

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

		// put data into redis
		if conn.Cmd("SET",sellerId, string(m.Data)).Err != nil {
			panic(fmt.Sprintf("Failed to perform REDIS operation"))
		}
	})

	// return lastscan data
	nc.Subscribe("dataservice.get.*.lastscan", func(m *nats.Msg) {
		fmt.Printf("Received a message. Subject: %s, Message: %s\n", m.Subject, string(m.Data))

		subjectSplit := strings.Split(m.Subject, ".")
		sellerId := subjectSplit[2]

		conn, err := redisPool.Get()
		if err != nil {
			panic(fmt.Sprintf("Error getting Redis connection: %s", err))
		}
		defer redisPool.Put(conn)

		// Get data from redis
		resp := conn.Cmd("GET",sellerId)
		if resp.Err != nil {
			panic(fmt.Sprintf("Failed to perform REDIS operation"))
		}

		body, err := resp.Bytes()
		if err != nil {
			fmt.Errorf("unable to read body: %s", err)
			nc.Publish(m.Reply, []byte("error"))
			return
		}

		nc.Publish(m.Reply, body)
	})

	<- killCh
}
