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
	streamer, err := newMessageStreamer()
	if err != nil {
		log.Fatal(err)
	}
	streamer.startStreaming()

	hb, err := newHeartbeater()
	if err != nil {
		log.Fatal(err)
	}

	fw := newFileWatcher(config.logFile, streamer.logStream)
	fw.startWatching()

	tokens := strings.Fields(config.serviceCommand)
	head := tokens[0]
	arguments := tokens[1:len(tokens)]

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
		logStderr(stderr)
		err = cmd.Wait()
		log.Printf("Service exited with %v", err)
		hb.stopBeating()
		time.Sleep(time.Second * 2)
	}
}

func logStderr(stderr io.ReadCloser) {
	go func() {
		for {
			r := bufio.NewReader(stderr)
			out, err := r.ReadString('\n')
			if err != nil {
				break
			}
			log.Printf("%s", out)
		}
	}()

}
