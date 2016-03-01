package logginghub

import (
	"bytes"
	"log"
	"net/http"
)

var (
	publish     chan []byte
	subscribe   chan []byte
	subscribers []subscriber
)

func startPubSubLoop() {
	publish = make(chan []byte, 1)
	subscribe = make(chan []byte, 1)
	subscribers = make([]subscriber, 10)

	go func() {
		for {
			select {
			case msg := <-publish:
				log.Printf("Publishing %s", string(msg[:]))
				for _, subscriber := range subscribers {
					subscriber.input(msg)
				}
			case msg := <-subscribe:
				log.Printf("Subscribing %s", string(msg[:]))
				subscribers = append(subscribers, subscriber{string(msg[:])})
			}
		}
	}()
}

type subscriber struct {
	Address string
}

func (s *subscriber) input(msg []byte) {
	go func() {
		_, _ = http.Post(s.Address, "application/text", bytes.NewReader(msg))
	}()
}
