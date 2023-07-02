package main

import (
	"encoding/json"
	"log"
	"sync"

	zmq "github.com/zeromq/goczmq"
)

const (
	endpoint = "tcp://127.0.0.1:5556"
)

type Payload struct {
	Index int
}

func subscriber() {
	// Create a sub channeler and connect it to the pub.
	sub, err := zmq.NewSub(endpoint, "payload")
	if err != nil {
		panic(err)
	}
	defer sub.Destroy()

	log.Println("subcriber connected")
	var payload Payload
	for {
		data, err := sub.RecvMessage()
		if err != nil {
			log.Panicln(err)
			continue
		}
		topic := data[0]
		buf := data[1]
		err = json.Unmarshal(buf, &payload)
		if err != nil {
			log.Panicln(err)
			continue
		}
		log.Printf("received: i=%d topic=%s", payload.Index, topic)
	}
}

func main() {
	log.Println("worker is starting")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		subscriber()
	}()

	log.Println("wait...")
	wg.Wait()
	log.Println("worker exit")
}
