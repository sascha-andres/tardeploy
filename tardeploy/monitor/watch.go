package monitor

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func Watch(directory string, batchInterval int, deployments chan<- string) {
	done := make(chan bool)
	defer close(done)

	log.Infof("Starting watcher with a batch interval of %ds", batchInterval)

	watcher, err := NewBatcher(time.Duration(batchInterval) * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		defer close(deployments)
		for {
			select {
			case event := <-watcher.Events:
				for key, _ := range event {
					deployments <- key
				}
			case err := <-watcher.Errors:
				log.Errorf("error:", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
