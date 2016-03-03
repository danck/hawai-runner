package runner

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

func init() {
	loadConfig()
}

var (
	streamer *messageStreamer
	fw       *fileWatcher
)

func Main() {
	//registerService()
	fw = newFileWatcher(config.logFile)
	fw.startWatching()

	streamer, _ = newMessageStreamer()
	streamer.startStreaming()

	initHeartbeat()

	tokens := strings.Fields(*serviceCommand)
	head := tokens[0]
	arguments := tokens[1:len(tokens)]

	//retries := 0
	// Run the guest service in an infinite loop
	for {
		startDelayedHeartbeat()
		cmd := exec.Command(head, arguments...)
		log.Println("Executing service command", *serviceCommand)
		out, err := cmd.Output()
		stopHeartbeat()
		if err != nil {
			log.Println("Error:", err.Error)
		}
		log.Println(string(out[:]))
		time.Sleep(time.Second * 2)
	}
}
