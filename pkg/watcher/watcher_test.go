package watcher

import (
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
)

func addPathToConsul(t *testing.T) {
	path := os.Getenv("PATH")
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("PATH", path+string(os.PathListSeparator)+dir)
}

func makeClientAndStartServer(t *testing.T) (*api.Client, *testutil.TestServer) {
	addPathToConsul(t)
	s, err := testutil.NewTestServerConfigT(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	conf := api.DefaultConfig()
	conf.Address = s.HTTPAddr
	c, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	return c, s
}

type entry struct {
	key   string
	kv    *api.KV
	delay time.Duration
	count int
}

func modify(e *entry) {
	for i := 0; i < e.count; i++ {
		e.kv.Put(&api.KVPair{Key: e.key, Value: []byte(strconv.Itoa(i))}, nil)
		time.Sleep(e.delay)
	}
}

func Test_ExecutesCallbackPerValueModification(t *testing.T) {
	c, s := makeClientAndStartServer(t)
	defer s.Stop()
	d := 10 * time.Millisecond
	total := 4

	w := NewWatcher("foo", d, c.KV())
	var wg sync.WaitGroup
	wg.Add(total)
	w.Start(func(v []byte) {
		wg.Done()
	}, nil)
	go func() {
		modify(&entry{
			key:   "foo/bar",
			kv:    c.KV(),
			delay: d,
			count: total / 2,
		})
	}()
	go func() {
		modify(&entry{
			key:   "foo/bar/1",
			kv:    c.KV(),
			delay: d,
			count: total / 2,
		})
	}()
	wg.Wait()
}

func Test_NotifiesOnKeyDelete(t *testing.T) {
	c, s := makeClientAndStartServer(t)
	defer s.Stop()

	kv := c.KV()

	p := &api.KVPair{Key: "baz", Value: []byte("1")}
	kv.Put(p, nil)

	w := NewWatcher("baz", 10*time.Millisecond, kv)
	var wg sync.WaitGroup
	wg.Add(1)
	w.Start(nil, func(k string) {

		wg.Done()
	})

	go func() {
		time.Sleep(10 * time.Millisecond)
		kv.Delete("baz", nil)
	}()
	wg.Wait()
}
