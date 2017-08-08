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

// taken from https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go
// and changed according to this use case

import (
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// Batcher batches file watch events in a given interval.
type Batcher struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}

	Events chan map[string]bool // Events are returned on this channel
}

// New creates and starts a Batcher with the given time interval.
func NewBatcher(interval time.Duration) (*Batcher, error) {
	watcher, err := fsnotify.NewWatcher()

	batcher := &Batcher{}
	batcher.Watcher = watcher
	batcher.interval = interval
	batcher.done = make(chan struct{}, 1)
	batcher.Events = make(chan map[string]bool, 1)

	if err == nil {
		go batcher.run()
	}

	return batcher, err
}

func (b *Batcher) run() {
	tick := time.Tick(b.interval)
	evs := make(map[string]bool, 0)
OuterLoop:
	for {
		select {
		case ev := <-b.Watcher.Events:
			if ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Rename == fsnotify.Rename {
				parts := strings.Split(ev.Name, "/")
				if _, ok := evs[parts[len(parts)-1]]; !ok {
					evs[parts[len(parts)-1]] = true
				}
			}
		case <-tick:
			log.Debug("Batch interval done")
			if len(evs) == 0 {
				continue
			}
			b.Events <- evs
			evs = make(map[string]bool, 0)
		case <-b.done:
			break OuterLoop
		}
	}
	close(b.done)
}

// Close stops the watching of the files.
func (b *Batcher) Close() {
	b.done <- struct{}{}
	b.Watcher.Close()
}
