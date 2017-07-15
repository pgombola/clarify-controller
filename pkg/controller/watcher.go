package controller

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
// sub-tree for changes. It polls at pollInterval.
func (w *Watcher) Start(f func(v []byte)) {
	go func() {
		ticker := time.NewTicker(w.pollInterval)
		for {
			select {
			case <-ticker.C:
				pairs, _, err := w.kv.List(w.Key, nil)
				if err != nil {
					// TODO: Actual logging -- gokit/log or logrus
					fmt.Printf("Error getting children: %v", err)
				}
				for _, p := range pairs {
					child, _ := w.children[p.Key]
					if child < p.ModifyIndex {
						f(p.Value)
					}
					w.children[p.Key] = p.ModifyIndex
				}
			case <-w.stop:
				ticker.Stop()
				// clear state
				w.children = make(map[string]uint64)
				return
			}
		}
	}()
}
