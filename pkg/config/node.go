package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Node struct {
	Hostname     string   `yaml:"Hostname"`
	NetInterface string   `yaml:"NetInterface"`
	Labels       []string `yaml:"Labels"`
}

type Nodes struct {
	nodes []*Node
}

func NewNodeConfig() *Nodes {
	return &Nodes{nodes: make([]*Node, 0)}
}

func (n *Nodes) Updated(v []byte) {
	node := &Node{}
	if err := yaml.Unmarshal(v, node); err != nil {
		fmt.Println(err)
	}
	n.nodes = append(n.nodes, node)
}

func (n *Nodes) Deleted(host string) {
	for i, node := range n.nodes {
		if node.Hostname == host {
			copy(n.nodes[i:], n.nodes[i+1:])
			n.nodes[len(n.nodes)-1] = &Node{}
			n.nodes = n.nodes[:len(n.nodes)-1]
		}
	}
}

func (n *Nodes) CurrentNode() (*Node, error) {
	host, err := os.Hostname()
	if err != nil {

	}
	for _, node := range n.nodes {
		if node.Hostname == host {
			return node, nil
		}
	}
	return nil, fmt.Errorf("unknown host: %s", host)
}
