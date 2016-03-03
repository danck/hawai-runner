package runner

import (
	"log"
)

type messageStreamer struct {
	loggingEndpoint string
	logStream       chan []byte
}

func newMessageStreamer() (*messageStreamer, error) {
	endpoint, err := resolveService("logging")
	if err != nil {
		//return nil, errors.New("Can't resolve logging endpoint", err.Error())
	}
	return &messageStreamer{
		loggingEndpoint: endpoint,
		logStream:       make(chan []byte, 1024),
	}, nil
}

func (ms *messageStreamer) startStreaming() {
	go ms.stream()
}

func (ms *messageStreamer) stream() {
	for {
		select {
		case msg := <-ms.logStream:
			log.Println(string(msg[:]))
		}
	}
}
