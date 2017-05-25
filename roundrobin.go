package liblb

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// RoundRobin implements round-robin alogrithm
type RoundRobin struct {
	i             int
	hosts         []string
	enableMetrics bool

	servedReqs *prometheus.CounterVec

	sync.Mutex
}

func NewRoundRobin(hosts ...string) *RoundRobin {
	return &RoundRobin{i: 0, hosts: hosts}
}

func (rb *RoundRobin) AddHost(host string) {
	rb.Lock()
	defer rb.Unlock()

	for _, h := range rb.hosts {
		if h == host {
			return
		}
	}
	rb.hosts = append(rb.hosts, host)
}

func (rb *RoundRobin) RemoveHost(host string) {
	rb.Lock()
	defer rb.Unlock()

	for i, h := range rb.hosts {
		if host == h {
			rb.hosts = append(rb.hosts[:i], rb.hosts[i+1:]...)
			break
		}
	}
}

func (rb *RoundRobin) EnableMetrics() error {
	rb.Lock()
	defer rb.Unlock()

	sreq := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "liblb_rb_requests_total",
		Help: "Number of requests served by RoundRobin balancer",
	}, []string{"host"})

	err := prometheus.Register(sreq)
	if err != nil {
		return err
	}
	rb.servedReqs = sreq

	rb.enableMetrics = true

	return nil
}

func (rb *RoundRobin) Balance() string {
	rb.Lock()
	defer rb.Unlock()

	if len(rb.hosts) == 0 {
		panic("no hosts")
	}

	host := rb.hosts[rb.i]
	rb.i++
	if rb.i >= len(rb.hosts) {
		rb.i = 0
	}

	if rb.enableMetrics {
		rb.servedReqs.WithLabelValues(host).Inc()
	}

	return host
}
