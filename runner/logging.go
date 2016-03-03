package runner

import (
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
)

type fileWatcher struct {
	file       string
	watcher    *fsnotify.Watcher
	offset     int64
	bufferSize int
	changed    chan bool
}

func (fw *fileWatcher) startWatching() {
	loggingService, err := getServiceAddress("logging")
	if err != nil {
		log.Println("No logging service available")
	}
	fw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go fw.notifier()
	go fw.streamer()

	err = fw.watcher.Add(*logFilePath)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func (fw *fileWatcher) notifier() {
	for {
		select {
		case event := <-fw.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				fw.changed <- true
			}
		case err := <-fw.watcher.Errors:
			log.Println("Error:", err)
		}
	}
}

func (fw *fileWatcher) streamer() {
	for {
		_ = <-fw.changed
		buffer := make([]byte, 1024)
		//open file
		file, err := os.Open(fw.file)
		if err != nil {
			log.Println("Error opening logfile:", err.Error())
			continue
		}
		defer file.Close()
		for {
			n, err := file.ReadAt(buffer, fw.offset)
			if err != io.EOF {
				log.Println("Error while reading logfile:", err.Error())
			}
			fw.offset += int64(n)
			logStream <- buffer
			if err == io.EOF {
				break
			}
		}
	}
}
