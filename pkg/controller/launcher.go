package controller

import "fmt"

type Launcher struct {
	Name string
	cmd  string
	args []string
}

// NewLauncher creates a launcher that will run
// applications and returns it.
func NewLauncher(name, cmd string) *Launcher {
	return &Launcher{cmd: cmd}
}

// Handle is called when a watcher notices a
// change in the keystore.
func (l *Launcher) Handle(v []byte) {

	fmt.Printf("Launching: %v\n", l.cmd)
}

// Checks if the current node can launch this
// application.
func canLaunch() bool {
	return false
}
