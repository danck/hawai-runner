package runner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	serviceName    string
	serviceAddress string
	logfile        string
)

const (
	registryAddress = "http://localhost:32000/service/"
)

func Main() {
	go registerAndKeepAlive()
	go observeAndPublish()
	select {}
}

func registerAndKeepAlive() {
	url := registryAddress + serviceName
	service := make(map[string]string)

	service["address"] = serviceAddress

	jsonStr, err := json.Marshal(service)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	id, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("Registered service with id %s", string(id[:]))

	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for _ = range ticker.C {

		}
	}()
}

func observeAndPublish() {
}
