package watcher

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type Watch struct {
	Key          string
	pollInterval time.Duration
	kv           *consulapi.KV
	children     map[string]uint64
	stop         chan bool
}

type Watcher interface {
	Updated([]byte)
	Deleted(string)
}

// NewWatcher creates a new watcher on the specified key.
func NewWatch(key string, pollInterval time.Duration, kv *consulapi.KV) *Watch {
	return &Watch{
		Key:          key,
		pollInterval: pollInterval,
		kv:           kv,
		children:     make(map[string]uint64),
		stop:         make(chan bool),
	}
}

// Stop exits from the watching go routine.
func (w *Watch) Stop() {
	w.stop <- true
}

// Start launches a go routine that watches a specified key
// sub-tree for changes. The modify indexes of keys that it
// finds are cached. If the modify index on the next polling
// interval are greater than the cached index, the providied
// function is called with the key's value.
func (w *Watch) Start(watcher Watcher) {
	go func() {
		ticker := time.NewTicker(w.pollInterval)
		for {
			select {
			case <-ticker.C:
				pairs, _, err := w.kv.List(w.Key, nil)
				if err != nil {
					// TODO: Actual logging -- gokit/log or logrus
					fmt.Printf("Error getting children for key %s: %v\n", w.Key, err)
				}
				updated := make(map[string]uint64)
				for _, p := range pairs {
					oldIndex, _ := w.children[p.Key]
					if oldIndex < p.ModifyIndex {
						watcher.Updated(p.Value)
					}
					updated[p.Key] = p.ModifyIndex
					delete(w.children, p.Key)
				}
				for k := range w.children {
					watcher.Deleted(k)
				}
				w.children = updated
			case <-w.stop:
				ticker.Stop()
				// clear cached indexes
				w.children = make(map[string]uint64)
				return
			}
		}
	}()
}
