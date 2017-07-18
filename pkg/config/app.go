package config

type App struct {
	Name     string `yaml:"Name"`
	Cmd      string
	Args     []string
	Count    int
	Liveness string
	Ports    map[string]string
}
