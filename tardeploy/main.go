package main

import (
	"fmt"
	"os"
	"os/signal"

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

	signals := make(chan os.Signal) // signal handling
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	deployments := make(chan string) // get deployment events
	go watch(deployments)

	done := make(chan bool)
	defer close(done)

loop:
	for {
		select {
		case s := <-signals: // handle signals
			log.Printf("Terminating program after receiving signal: %v", s)
			break loop
		case deployment, ok := <-deployments: // handle deployment events
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
	log.Printf("Checking %s", path.Join(configuration.Directories.TarballDirectory, deployment))
	ok, err := exists(path.Join(configuration.Directories.TarballDirectory, deployment))
	if ok && err == nil {
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
