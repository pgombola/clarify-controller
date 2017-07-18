package main

import (
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/pgombola/consul-watch-test/pkg/watcher"
)

func main() {
	// Parse flags -- Will git flag work or do we need viper/cobra?

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
	appCfgWatch := watcher.NewWatch("/registry/config/apps", 1*time.Minute, client.KV())

	nodeCfgWatch := watcher.NewWatch("/registry/config/nodes", 1*time.Minute, client.KV())

	sysCfgWatch := watcher.NewWatch("/registry/config/system", 1*time.Minute, client.KV())

	// Create app watcher
	appWatch := watcher.NewWatch("/registry/app", 15*time.Second, client.KV())

	// Create lifecycle watcher -- or is this part of /registry/config/nodes?

}
