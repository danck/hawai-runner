package logginghub

import (
	"io/ioutil"
	"log"
	"net/http"
)

// pubHandler receives lines of text
func pubHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit pubHandler")
	logData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		r.Body.Close()
		return
	}
	publish <- logData
}

// subHandler registers subscibers
func subHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Hit subHandler")
	subscriberAddress, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		r.Body.Close()
		return
	}
	subscribe <- subscriberAddress
}
