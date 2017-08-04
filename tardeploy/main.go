package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/log"

	"path"

	"github.com/google/gops/agent"
	"github.com/sascha-andres/tardeploy"
)

var (
	configuration *tardeploy.Configuration
)

func main() {
	// gops
	if err := agent.Listen(nil); err != nil {
		log.Fatalln(err)
	}

	configuration = config() // load configuration and validate

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	deployments := make(chan string)

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
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(configuration.Directories.TarballDirectory)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	defer close(done)

loop:
	for {
		select {
		case s := <-signals:
			log.Printf("Terminating program after receiving signal: %v", s)
			break loop
		case deployment, ok := <-deployments:
			if !ok {
				break loop
			}
			if "" != deployment {
				go handleChange(deployment)
			}
		}
	}

	<-done
}

func handleChange(deployment string) {
	ok, err := exists(path.Join(configuration.Directories.TarballDirectory, deployment))
	if ok && err != nil {
		if err := configuration.SetupApplication(deployment); err != nil {
			log.Warnln(fmt.Sprintf("Error deploying application %s: %#v", deployment, err))
		}
	} else {
		if err != nil {
			log.Warnln(fmt.Sprintf("Error deploying application %s: %#v", deployment, err))
		} else {
			if err := configuration.RemoveApplication(deployment); err != nil {
				log.Warnln(fmt.Sprintf("Error removing application %s: %#v", deployment, err))
			}
		}
	}
}
