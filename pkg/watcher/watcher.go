package watcher

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type Watcher struct {
	Key          string
	pollInterval time.Duration
	kv           *consulapi.KV
	children     map[string]uint64
	stop         chan bool
}

// NewWatcher creates a new watcher on the specified key.
func NewWatcher(key string, pollInterval time.Duration, kv *consulapi.KV) *Watcher {
	return &Watcher{
		Key:          key,
		pollInterval: pollInterval,
		kv:           kv,
		children:     make(map[string]uint64),
		stop:         make(chan bool),
	}
}

// Stop exits from the watching go routine.
func (w *Watcher) Stop() {
	w.stop <- true
}

// Start launches a go routine that watches a specified key
// sub-tree for changes. The modify indexes of keys that it
// finds are cached. If the modify index on the next polling
// interval are greater than the cached index, the providied
// function is called with the key's value.
func (w *Watcher) Start(mod func(v []byte), del func(k string)) {
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
					if oldIndex < p.ModifyIndex && mod != nil {
						mod(p.Value)
					}
					updated[p.Key] = p.ModifyIndex
					delete(w.children, p.Key)
				}
				for k := range w.children {
					if del != nil {
						del(k)
					}
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
