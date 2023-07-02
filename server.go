package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	czmq "github.com/zeromq/goczmq"
)

const (
	endpoint = "tcp://127.0.0.1:5556"
)

type Payload struct {
	Index int
}

func publisher() {
	topic := []byte("payload")
	// Create a pub channeler and bind it to port 5555.
	// A channeler provides a thread safe channel interface
	// to a *Sock
	pub, err := czmq.NewPub(endpoint)

	if err != nil {
		panic(err)
	}
	defer pub.Destroy()
	// pub.Bind(endpoint)
	// if err != nil {
	// 	panic(err)
	// }
	log.Printf("publisher bound")
	// pub.SetSubscribe("/payload")

	time.Sleep(1 * time.Second)
	// preload := []byte("welcome")
	// pub.SendMessage([][]byte{topic, preload})
	// pub.SendMessage([][]byte{topic, preload})

	for i := 0; i < 100; i++ {
		payload := Payload{Index: i}
		msg, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		err = pub.SendMessage([][]byte{topic, msg})
		if err != nil {
			panic(err)
		}
		log.Printf("publish: i=%d topic=%s", i, topic)
		time.Sleep(1 * time.Millisecond)
	}

}

func main() {
	log.Println("server is starting")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		publisher()
	}()

	wg.Wait()
	log.Println("server is done")
}
