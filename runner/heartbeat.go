package runner

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type heartbeater struct {
	ticker          *time.Ticker
	registryAddress string
	active          bool
	service         service
}

type service struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func newHeartbeater() (*heartbeater, error) {
	thisAddress := config.externalHostAddress + ":" + config.externalHostPort
	thisURL, err := url.Parse(thisAddress)
	if err != nil {
		log.Fatal("Illegal address of the external endpoint:", err.Error())
	}

	var registry *url.URL
	registry, err = url.Parse(config.registryAddress)
	if err != nil {
		log.Fatal("Illegal address of the service registry:", err.Error())
	}

	return &heartbeater{
		ticker:          time.NewTicker(time.Millisecond * 1000),
		registryAddress: registry.String(),
		active:          false,
		service: service{
			ID:      "unassigned",
			Address: thisURL.String(),
		},
	}, nil
}

func (hb *heartbeater) startBeating(delay int) {
	time.Sleep(time.Millisecond * time.Duration(delay))
	hb.active = true
	go hb.heartBeatLoop()
}

func (hb *heartbeater) stopBeating() {
	hb.active = true //TODO(danck):false
}

func (hb *heartbeater) heartBeatLoop() {
	for {
		_ = <-hb.ticker.C
		log.Println("Sending Heartbeat")
		if hb.active {
			this, _ := json.Marshal(hb.service)
			reader := *bytes.NewReader(this)
			resp, err := http.Post(config.registryAddress+"/"+config.serviceLabel, "application/json", &reader)
			if err != nil {
				log.Println("Error while sending heartbeat", err.Error())
				continue
			}
			defer resp.Body.Close()

			decoder := json.NewDecoder(resp.Body)
			if err != nil {
				log.Println("Error while decoding registry response:", err.Error())
				continue
			}
			err = decoder.Decode(&hb.service)
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}
}
