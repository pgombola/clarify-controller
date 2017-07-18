package launcher

import (
	"fmt"

	"github.com/pgombola/consul-watch-test/pkg/config"
)

type Launcher struct {
	Name string
}

// NewLauncher creates a launcher that will run
// applications and returns it.
func NewLauncher(name string) *Launcher {
	return &Launcher{Name: name}
}

// Handle is called when a watcher notices a
// change in the keystore.
func (l *Launcher) Handle(v []byte) {
	app := config.ParseAppConfig(v)
	fmt.Printf("Launching: %v\n", l.cmd)
}

// Checks if the current node can launch this
// application.
func canLaunch(n *config.Node) bool {
	for _, l := range n.Labels {
		if 
	}
	return false
}
