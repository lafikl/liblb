// R2 is a concurrency-safe Round-Robin Balancer.
// Which also supports Weighted Round-Robin.
//
// Round-Robin is a simple and well known algorithm for load balancing.
//
// Here's a simple implmentation of it, if you're not already familiar with it
//
// https://play.golang.org/p/XCMAtKGCaE
package r2

import (
	"sync"

	"github.com/lafikl/liblb"
)

type R2 struct {
	i     int
	hosts []string

	sync.Mutex
}

func New(hosts ...string) *R2 {
	return &R2{i: 0, hosts: hosts}
}

// Adds a host to the list of hosts, with the weight of the host being 1.
func (rb *R2) Add(host string) {
	rb.Lock()
	defer rb.Unlock()

	for _, h := range rb.hosts {
		if h == host {
			return
		}
	}
	rb.hosts = append(rb.hosts, host)
}

// Weight increases the percentage of requests that get sent to the host
// Which can be calculated as `weight/(total_weights+weight)`.
func (rb *R2) AddWeight(host string, weight int) {
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

// Check if host already exist
func (rb *R2) Exists(host string) bool {
	rb.Lock()
	defer rb.Unlock()

	for _, h := range rb.hosts {
		if h == host {
			return true
		}
	}

	return false
}

func (rb *R2) Remove(host string) {
	rb.Lock()
	defer rb.Unlock()

	for i, h := range rb.hosts {
		if host == h {
			rb.hosts = append(rb.hosts[:i], rb.hosts[i+1:]...)
		}
	}
}

func (rb *R2) Balance() (string, error) {
	rb.Lock()
	defer rb.Unlock()

	if len(rb.hosts) == 0 {
		return "", liblb.ErrNoHost
	}

	host := rb.hosts[rb.i%len(rb.hosts)]
	rb.i++

	return host, nil
}
