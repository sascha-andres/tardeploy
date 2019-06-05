// Copyright © 2017 Sascha Andres <sascha.andres@outlook.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"path"

	"github.com/google/gops/agent"
	"github.com/sascha-andres/tardeploy"
	"github.com/sascha-andres/tardeploy/tardeploy/monitor"
)

var (
	configuration *tardeploy.Configuration
)

func setupLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	default:
		log.SetLevel(log.InfoLevel)
		break
	}
	log.Infof("Log level set to %s", level)
}

func main() {
	log.WithField("version", version).Info("tardeploy")

	// gops
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalln(err)
	}

	configuration = config() // load configuration and validate
	setupLogLevel(configuration.Application.LogLevel)

	signals := make(chan os.Signal) // signal handling
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	deployments := make(chan string) // get deployment events
	go monitor.Watch(configuration.Directories.TarballDirectory, configuration.Application.BatchInterval, deployments)

	done := make(chan bool)
	defer close(done)

loop:
	for {
		select {
		case s := <-signals: // handle signals
			log.Infof("Terminating program after receiving signal: %v", s)
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
	log.Infof("Checking %s", path.Join(configuration.Directories.TarballDirectory, deployment))
	ok, err := exists(path.Join(configuration.Directories.TarballDirectory, deployment))
	if ok && err == nil {
		if err := configuration.SetupApplication(deployment); err != nil {
			log.Warnf("Error deploying application %s: %#v", deployment, err)
		}
	} else {
		if err != nil {
			log.Warnf("Error deploying application %s: %#v", deployment, err)
		} else {
			if err := configuration.RemoveApplication(deployment); err != nil {
				log.Warnf("Error removing application %s: %#v", deployment, err)
			}
		}
	}
}
