package logginghub

import (
	"flag"
	"log"
	"net/http"
)

var (
	listenAddress = flag.String(
		"listen-addr",
		":20000",
		"Address to listen for incoming log posts")
)

func Main() {
	flag.Parse()
	router := http.NewServeMux()
	router.HandleFunc("/pub", pubHandler)
	router.HandleFunc("/sub", subHandler)

	startPubSubLoop()

	log.Printf("Starting to listen on %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, router))
}
