package config

type Cluster struct {
	Share   string `yaml:"Share"`
	Version map[string]string
}
