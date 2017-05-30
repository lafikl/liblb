package liblb

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// R2 implements round-robin alogrithm
type R2 struct {
	i             int
	hosts         []string
	enableMetrics bool

	servedReqs *prometheus.CounterVec

	sync.Mutex
}

func NewR2(hosts ...string) *R2 {
	return &R2{i: 0, hosts: hosts}
}

func (rb *R2) AddHost(host string) {
	rb.Lock()
	defer rb.Unlock()

	for _, h := range rb.hosts {
		if h == host {
			return
		}
	}
	rb.hosts = append(rb.hosts, host)
}

func (rb *R2) AddHostWithWeight(host string, weight int) {
	rb.Lock()
	defer rb.Unlock()

	for _, h := range rb.hosts {
		if h == host {
			return
		}
	}

	for i := 0; i < weight; i++ {
		rb.hosts = append(rb.hosts, host)
	}

}

func (rb *R2) HostExists(host string) bool {
	rb.Lock()
	defer rb.Unlock()

	for _, h := range rb.hosts {
		if h == host {
			return true
		}
	}

	return false
}

func (rb *R2) RemoveHost(host string) {
	rb.Lock()
	defer rb.Unlock()

	for i, h := range rb.hosts {
		if host == h {
			rb.hosts = append(rb.hosts[:i], rb.hosts[i+1:]...)
		}
	}
}

func (rb *R2) EnableMetrics() error {
	rb.Lock()
	defer rb.Unlock()

	sreq := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "liblb_r2_requests_total",
		Help: "Number of requests served by R2 balancer",
	}, []string{"host"})

	err := prometheus.Register(sreq)
	if err != nil {
		return err
	}
	rb.servedReqs = sreq

	rb.enableMetrics = true

	return nil
}

func (rb *R2) Balance() (string, error) {
	rb.Lock()
	defer rb.Unlock()

	if len(rb.hosts) == 0 {
		return "", ErrNoHost
	}

	host := rb.hosts[rb.i]
	rb.i++
	if rb.i >= len(rb.hosts) {
		rb.i = 0
	}

	if rb.enableMetrics {
		rb.servedReqs.WithLabelValues(host).Inc()
	}

	return host, nil
}
