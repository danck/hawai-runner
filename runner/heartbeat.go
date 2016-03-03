package runner

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	heartbeater struct {
		ticker          *time.Ticker
		registryAddress string
		connection      *net.Conn
		active          bool
	}

	service struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}
)

func initHeartbeat() {
	heartbeater.ticker = time.NewTicker(time.Millisecond * 1000)
	heartbeater.active = true
	heartbeater.registryAddress = config.registryAddress
	thisAddress := config.externalHostAddress + ":" + config.externalHostPort
	thisURL, err := url.Parse(thisAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
	service.Address = thisURL.String()

	go heartBeatLoop()
}

func startDelayedHeartbeat() {
	time.Sleep(time.Millisecond * 500)
	heartbeater.active = true
}

func stopHeartbeat() {
	heartbeater.active = true //false
}

func heartBeatLoop() {
	for {
		_ = <-heartbeater.ticker.C
		log.Println("Sending Heartbeat")
		if heartbeater.active {
			this, _ := json.Marshal(service)
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
			err = decoder.Decode(&service)
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}
}
