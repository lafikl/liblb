package p2c

import (
	"hash/fnv"
	"math/rand"
	"sync"
	"time"

	"github.com/lafikl/liblb"
	"github.com/lafikl/liblb/murmur"
	"github.com/prometheus/client_golang/prometheus"
)

type host struct {
	name string
	load uint64
}

// P2C will distribute the traffic by choosing two hosts either via hashing or randomly
// and then pick the least loaded of the two.
// It gaurantees that the max load of a server is ln(ln(n)),
// where n is the number of servers.
//
// ## Load in P2C:
//
// For more on P2C http://www.eecs.harvard.edu/~michaelm/postscripts/handbook2001.pdf
//
type P2C struct {
	hosts   []*host
	rndm    *rand.Rand
	loadMap map[string]*host

	enableMetrics bool
	servedReqs    *prometheus.CounterVec
	pservedReqs   *prometheus.CounterVec

	sync.Mutex
}

// New returns a new instance of RandomTwoBalancer
func New(hosts ...string) *P2C {
	p := &P2C{
		hosts:   []*host{},
		loadMap: map[string]*host{},
		rndm:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, h := range hosts {
		p.Add(h)
	}
	return p
}

// Register liblb_p2c_requests_total metric in prometheus
func (p *P2C) EnableMetrics() error {
	p.Lock()
	defer p.Unlock()

	sreq := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "liblb_p2c_requests_total",
		Help: "Number of requests served by P2C balancer",
	}, []string{"host"})

	err := prometheus.Register(sreq)
	if err != nil {
		return err
	}
	p.servedReqs = sreq

	psreq := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "liblb_pp2c_requests_total",
		Help: "Number of requests served by Partial Key balancer",
	}, []string{"host"})

	err = prometheus.Register(psreq)
	if err != nil {
		return err
	}
	p.pservedReqs = psreq

	p.enableMetrics = true
	return nil
}

func (p *P2C) Add(hostName string) {
	p.Lock()
	defer p.Unlock()

	h := &host{name: hostName, load: 0}
	p.hosts = append(p.hosts, h)
	p.loadMap[hostName] = h
}

func (p *P2C) Remove(host string) {
	p.Lock()
	defer p.Unlock()

	_, ok := p.loadMap[host]
	if !ok {
		return
	}

	delete(p.loadMap, host)

	for i, v := range p.hosts {
		if v.name == host {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
		}
	}
}

func (p *P2C) hash(key string) (string, string) {
	h := fnv.New32()
	h.Write([]byte(key))

	n1 := p.hosts[int(h.Sum32())%len(p.hosts)].name
	n2 := p.hosts[int(murmur.Murmur3([]byte(key)))%len(p.hosts)].name

	return n1, n2

}

// Balance picks two servers either randomly (if no key supplied), or via
// hashing (PKG) if given a key, then it returns the least loaded one between the two.
//
// Partial Key Grouping (PKG) is great for skewed data workloads, which also needs to be
// determinstic in the way of choosing which servers to send requests too.
// https://arxiv.org/pdf/1510.07623.pdf
// the maximum load of a server in PKG at anytime is:
// `max_load-avg_load`
func (p *P2C) Balance(key string) (string, error) {
	p.Lock()
	defer p.Unlock()

	if len(p.hosts) == 0 {
		return "", liblb.ErrNoHost
	}

	// chosen host
	var host string

	var n1, n2 string

	if len(key) > 0 {
		n1, n2 = p.hash(key)
	} else {
		n1 = p.hosts[p.rndm.Intn(len(p.hosts))].name
		n2 = p.hosts[p.rndm.Intn(len(p.hosts))].name
	}

	host = n2

	if p.loadMap[n1].load <= p.loadMap[n2].load {
		host = n1
	}

	if p.enableMetrics {
		if len(key) > 0 {
			p.pservedReqs.WithLabelValues(host).Inc()
		} else {
			p.servedReqs.WithLabelValues(host).Inc()
		}
	}

	p.loadMap[host].load++
	return host, nil
}

// Decrments the load of the host (if found) by 1
func (p *P2C) Done(host string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return
	}
	if h.load > 0 {
		h.load--
	}
}

// UpdateLoad updates the load of a host
func (p *P2C) UpdateLoad(host string, load uint64) error {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return liblb.ErrNoHost
	}
	h.load = load
	return nil
}

// Returns the current load of the server,
// or it returns liblb.ErrNoHost if the host doesn't exist.
func (p *P2C) GetLoad(host string) (load uint64, err error) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return 0, liblb.ErrNoHost
	}
	return h.load, nil
}
