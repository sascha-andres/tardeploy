package monitor

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/log"
)

func Watch(directory string, deployments chan<- string) {
	done := make(chan bool)
	defer close(done)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		defer close(deployments)
		for {
			select {
			case event := <-watcher.Events:
				// event.Op&fsnotify.Create == fsnotify.Create ||
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Write == fsnotify.Write {
					parts := strings.Split(event.Name, "/")
					deployments <- parts[len(parts)-1]
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
