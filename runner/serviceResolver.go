package runner

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type services struct {
	list []string `json:"services"`
}

// getServiceAddress queries the service registry for the available addresses
// of a given label
func resolveService(label string) (string, error) {
	serviceURL, err := url.Parse(*registryAddress + "/" + label)
	if err != nil {
		log.Fatal("Error while getting an service address", err.Error())
	}
	resp, err := http.Get(serviceURL.String())
	if err != nil {
		log.Fatal("Error while querying the service registry", err.Error())
	}
	defer resp.Body.Close()

	var serviceList services
	err = json.NewDecoder(resp.Body).Decode(&serviceList)
	if err != nil {
		log.Fatal("Error while decoding response from service registry", err.Error())
	}
	if len(serviceList.list) == 0 {
		return "", errors.New("No services available")
	}
	chosenAddress := serviceList.list[0]

	return chosenAddress, nil
}
