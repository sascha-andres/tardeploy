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

package monitor

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Watch starts looking for files placed in directory
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
				for key := range event {
					deployments <- key
				}
			case err := <-watcher.Errors:
				log.Errorf("Error: %s", err.Error())
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
