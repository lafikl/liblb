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
	"github.com/lafikl/liblb"
	"github.com/tevino/abool"
)

type R2 struct {
	i     int
	hosts []string

	lock *abool.AtomicBool
}

func New(hosts ...string) *R2 {
	return &R2{i: 0, hosts: hosts, lock: abool.New()}
}

// Adds a host to the list of hosts, with the weight of the host being 1.
func (rb *R2) Add(host string) {
	rb.lock.Set()
	defer rb.lock.UnSet()

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
	rb.lock.Set()
	defer rb.lock.UnSet()

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
	rb.lock.Set()
	defer rb.lock.UnSet()

	for _, h := range rb.hosts {
		if h == host {
			return true
		}
	}

	return false
}

func (rb *R2) Remove(host string) {
	rb.lock.Set()
	defer rb.lock.UnSet()

	for i, h := range rb.hosts {
		if host == h {
			rb.hosts = append(rb.hosts[:i], rb.hosts[i+1:]...)
		}
	}
}

func (rb *R2) Balance() (string, error) {
	rb.lock.Set()
	defer rb.lock.UnSet()

	if len(rb.hosts) == 0 {
		return "", liblb.ErrNoHost
	}

	host := rb.hosts[rb.i%len(rb.hosts)]
	rb.i++

	return host, nil
}
