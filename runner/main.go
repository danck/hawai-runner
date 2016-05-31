package runner

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

func init() {
	loadConfig()
}

func Main() {
	hb, err := newHeartbeater()
	if err != nil {
		log.Fatal(err)
	}
	hb.startBeating(0)
	time.Sleep(1 * time.Second)

	streamer, err := newMessageStreamer()
	if err != nil {
		log.Fatal(err)
	}
	streamer.startStreaming()

	fw := newFileWatcher(config.logFile, streamer.logStream)
	fw.startWatching()

	tokens := strings.Fields(config.serviceCommand)
	head := tokens[0]
	arguments := tokens[1:len(tokens)]

	hb.stopBeating()
	// Run the guest service in an infinite loop
	for {
		hb.startBeating(1000)
		cmd := exec.Command(head, arguments...)
		stderr, err := cmd.StderrPipe()
		log.Println("Executing service command", config.serviceCommand)
		err = cmd.Start()
		if err != nil {
			log.Println("Error:", err.Error)
		}
		logStderr(stderr, streamer)
		err = cmd.Wait()
		log.Printf("Service exited with %v", err)
		hb.stopBeating()
		time.Sleep(time.Second * 2)
	}
}

func logStderr(stderr io.ReadCloser, ms *messageStreamer) {
	r := bufio.NewReader(stderr)
	go func() {
		for {
			out, err := r.ReadString('\n')
			if err != nil {
				errMsg := "Error while reading stderr: " + err.Error()
				ms.logStream <- []byte(errMsg)
				log.Println(errMsg)
				break
			}
			ms.logStream <- []byte(out)
			log.Printf("%s", out)
		}
	}()
}
