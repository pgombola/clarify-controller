package config

import "testing"

func Test_AddNode(t *testing.T) {
	n := NewNodeConfig()
	n.Updated([]byte("Hostname: foo1"))

	if n.nodes[0].Hostname != "foo1" {
		t.Fatal("Expected \"foo\", got \"" + n.nodes[0].Hostname + "\"")
	}
}

func Test_Add2Nodes(t *testing.T) {
	n := NewNodeConfig()
	n.Updated([]byte("Hostname: foo1"))
	n.Updated([]byte("Hostname: foo2"))

	if len(n.nodes) != 2 {
		t.Fatalf("Expected length of 2, got length of %d\n", len(n.nodes))
	}
}

func Test_DeleteFirstNode(t *testing.T) {
	nc := NewNodeConfig()
	nc.Updated([]byte("Hostname: foo1"))
	nc.Updated([]byte("Hostname: foo2"))

	nc.Deleted("foo1")
	if len(nc.nodes) != 1 {
		t.Fatalf("Expected length of 1, got length of %d", len(nc.nodes))
	}
}

func Test_DeleteLastNode(t *testing.T) {
	nc := NewNodeConfig()
	nc.Updated([]byte("Hostname: foo1"))
	nc.Updated([]byte("Hostname: foo2"))

	nc.Deleted("foo2")
	if len(nc.nodes) != 1 {
		t.Fatalf("Expected length of 1, got length of %d", len(nc.nodes))
	}
}

func Test_Delete3Nodes(t *testing.T) {
	nc := NewNodeConfig()
	nc.Updated([]byte("Hostname: foo1"))
	nc.Updated([]byte("Hostname: foo2"))
	nc.Updated([]byte("Hostname: foo3"))
	nc.Updated([]byte("Hostname: foo4"))
	nc.Updated([]byte("Hostname: foo5"))
	nc.Updated([]byte("Hostname: foo6"))
	nc.Updated([]byte("Hostname: foo7"))

	nc.Deleted("foo3")
	nc.Deleted("foo4")
	nc.Deleted("foo5")

	if len(nc.nodes) != 4 {
		t.Fatalf("Expected length of 4, got length of %d", len(nc.nodes))
	}
}

func Test_DeleteEmptyNodes(t *testing.T) {
	nc := NewNodeConfig()
	nc.Deleted("foo")
	if len(nc.nodes) > 0 {
		t.Fatalf("Expected length of 0, got length of %d", len(nc.nodes))
	}
}
