// Bounded is Consistent hashing with bounded loads.
// It acheives that by adding a capacity counter on every host,
// and when a host gets picked it, checks its capacity to see if it's below
// the Average Load per Host.
//
// All opertaions in bounded are concurrency-safe.
//
// Average Load Per Host is defined as follows:
//
// (totalLoad/number_of_hosts)*imbalance_constant
//
// totalLoad = sum of all hosts load
//
// load = the number of active requests
//
// imbalance_constant = is the imbalance constant, which is 1.25 in our case
//
// it bounds the load imabalnce to be at most 25% more than (totalLoad/number_of_hosts).
//
//
// For more info:
// https://medium.com/vimeo-engineering-blog/improving-load-balancing-with-a-new-consistent-hashing-algorithm-9f1bd75709ed
//
// https://research.googleblog.com/2017/04/consistent-hashing-with-bounded-loads.html
package bounded

import (
	"github.com/lafikl/consistent"
	"github.com/lafikl/liblb"
)

type bhost struct {
	load   uint64
	weight int
}

type Bounded struct {
	ch *consistent.Consistent
}

func New(hosts ...string) *Bounded {
	c := &Bounded{
		ch: consistent.New(),
	}
	for _, h := range hosts {
		c.Add(h)
	}
	return c
}

func (b *Bounded) Add(host string) {
	b.ch.Add(host)
}

func (b *Bounded) Remove(host string) {
	b.ch.Remove(host)
}

// err can be liblb.ErrNoHost if there's no added hosts.
func (b *Bounded) Balance(key string) (host string, err error) {
	if len(b.ch.Hosts()) == 0 {
		return "", liblb.ErrNoHost
	}

	host, err = b.ch.GetLeast(key)
	return
}

// It should be called once a request is assigned to a host,
// obtained from b.Balance.
func (b *Bounded) Inc(host string) {
	b.ch.Inc(host)
}

// should be called when an assigned request to host is finished.
func (b *Bounded) Done(host string) {
	b.ch.Done(host)
}

func (b *Bounded) Loads() map[string]int64 {
	return b.ch.GetLoads()
}

// Max load of a host is (Average Load Per Host*1.25)
func (b *Bounded) MaxLoad() int64 {
	return b.ch.MaxLoad()
}
