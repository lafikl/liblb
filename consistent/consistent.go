package consistent

import (
	"sync"

	"github.com/lafikl/liblb"
	"github.com/prometheus/client_golang/prometheus"

	"stathat.com/c/consistent"
)

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

func (c *Consistent) AddHost(host string) {
	c.ch.Add(host)
}

func (c *Consistent) RemoveHost(host string) {
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
