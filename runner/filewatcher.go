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
	outStream  chan []byte
}

func newFileWatcher(file string, out chan []byte) *fileWatcher {
	return &fileWatcher{
		file:       file,
		offset:     int64(0),
		bufferSize: 4096,
		changed:    make(chan bool, 16),
		outStream:  out,
	}
}

func (fw *fileWatcher) startWatching() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	fw.watcher = watcher
	go fw.notifier()
	go fw.streamer()

	err = fw.watcher.Add(fw.file)
	if err != nil {
		log.Printf("Logfile not found: %v", err)
		log.Printf("Creating empty file at: %s", fw.file)
		_, err := os.Create(fw.file)
		if err != nil {
			log.Fatal("Failed to create empty logfile", err.Error())
		}
	}
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
			fw.outStream <- buffer
			if err == io.EOF {
				break
			}
		}
	}
}
