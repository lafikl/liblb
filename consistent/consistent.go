package consistent

import (
	"sync"

	"github.com/lafikl/liblb"
	"github.com/prometheus/client_golang/prometheus"

	"stathat.com/c/consistent"
)

// Consistent uses consistent hashing algorithm to assign work to hosts.
// Its best for the cases when you need affinty, and your hosts come and go.
// When removing a host it gaurantees that only 1/n of items gets reshuffled
// where n is the number of servers.
//
// One of the issues with Consistent Hashing is load imbalance
// when you hot keys that goes to a single server,
// it's mitigated by using virtual nodes,
// which basically means when adding a host we add n - 20 in our case - copies.
//
// Beware that if Consistent Hashing doesn't provide,
// an upper bound for the load of a host.
// If you need such gaurantees see package bounded.
//
// https://en.wikipedia.org/wiki/Consistent_hashing
type Consistent struct {
	ch *consistent.Consistent

	enableMetrics bool
	servedReqs    *prometheus.CounterVec
	sync.RWMutex
}

func New(hosts ...string) *Consistent {
	c := &Consistent{ch: consistent.New()}
	for _, h := range hosts {
		c.ch.Add(h)
	}
	return c
}

func (c *Consistent) EnableMetrics() error {
	c.Lock()
	defer c.Unlock()

	sreq := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "liblb_consistent_requests_total",
		Help: "Number of requests served by Consistent balancer",
	}, []string{"host"})

	err := prometheus.Register(sreq)
	if err != nil {
		return err
	}
	c.servedReqs = sreq
	c.enableMetrics = true
	return nil
}

func (c *Consistent) Add(host string) {
	c.ch.Add(host)
}

func (c *Consistent) Remove(host string) {
	c.ch.Remove(host)
}

func (h *Consistent) Balance(key string) (host string, err error) {
	host, err = h.ch.Get(key)
	if err != nil {
		return "", liblb.ErrNoHost
	}
	if h.enableMetrics {
		h.servedReqs.WithLabelValues(host).Inc()
	}
	return
}
