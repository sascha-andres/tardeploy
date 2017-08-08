package monitor

// taken from https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go

import (
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
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
			if ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Write == fsnotify.Write {
				parts := strings.Split(ev.Name, "/")
				if _, ok := evs[parts[len(parts)-1]]; !ok {
					evs[parts[len(parts)-1]] = true
				}
			}
		case <-tick:
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
