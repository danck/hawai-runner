package runner

import (
	"bytes"
	"errors"
	"log"
	"net/http"
)

const (
	apiPrefix  = "/pub"
	apiChannel = "/logging"
)

type messageStreamer struct {
	loggingEndpoint string
	logStream       chan []byte
}

func newMessageStreamer() (*messageStreamer, error) {
	endpoint, err := resolveService("logging")
	if err != nil {
		return nil, errors.New("Can't resolve logging endpoint" + err.Error())
	}
	endpoint = endpoint + apiPrefix + apiChannel
	log.Println("Logging endpoint:", endpoint)
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
			body := bytes.NewReader(msg)
			resp, err := http.Post(ms.loggingEndpoint, "application/text", body)
			if err != nil {
				log.Printf("Failed to post message. Reason: %s", err.Error())
			}
			defer resp.Body.Close()
		}
	}
}
