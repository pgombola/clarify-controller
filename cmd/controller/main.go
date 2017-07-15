package main

import (
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/pgombola/consul-watch-test/pkg/controller"
)

func main() {
	// Parse flags -- Will flag work or do we need viper/cobra?

	// Run updater

	// Run consul

	// Create consul client
	// We'll need to change from DefaultConfig when consul port becomes configurable
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		log.Fatalf("couldn't connect to consul: %v", err)
	}

	// Create config watchers
	// TODO: Pass logger in
	appwatch := controller.NewWatcher("/registry/config/apps", 1*time.Minute, client.KV())
	appwatch.Start(func(v []byte) {
		// if running, restart
		// Where do we check to see if the current node can launch the application?
		// if not running, create launcher based on unmarshalled v
	})
	controller.NewWatcher("/registry/config/nodes", 1*time.Minute, client.KV())
	controller.NewWatcher("/registry/config/system", 1*time.Minute, client.KV())

	// Create app watcher
	controller.NewWatcher("/registry/app", 15*time.Second, client.KV())

	// Create lifecycle watcher -- or is this part of /registry/config/nodes?

}
