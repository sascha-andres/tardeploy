package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"time"

	"path"

	"github.com/google/gops/agent"
	"github.com/sascha-andres/tardeploy"
	"github.com/sascha-andres/tardeploy/handler"
)

var (
	configuration tardeploy.Configuration
)

func main() {
	// gops
	if err := agent.Listen(nil); err != nil {
		log.Fatal(err)
	}

	config() // load configuration and validate

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	deployments := make(chan string)
	go func() {
		// TODO setup watcher [ get changes file, pass to handler ]
		defer close(deployments)
		time.Sleep(5 * time.Second)
		deployments <- "hallo"
	}()

loop:
	for {
		select {
		case s := <-signals:
			log.Printf("Terminating program after receiving signal: %v\n", s)
			break loop
		case deployment := <-deployments:
			go handleChange(deployment)
		}
	}
}
func handleChange(deployment string) {
	ok, err := exists(path.Join(configuration.TarballDirectory, deployment))
	if ok && err != nil {
		if err := handler.SetupApplication(deployment, configuration); err != nil {
			log.Println(fmt.Sprintf("Error deploying application %s: %#v", deployment, err))
		}
	} else {
		if err != nil {
			log.Println(fmt.Sprintf("Error deploying application %s: %#v", deployment, err))
		} else {
			if err := handler.RemoveApplication(deployment, configuration); err != nil {
				log.Println(fmt.Sprintf("Error removing application %s: %#v", deployment, err))
			}
		}
	}
}
