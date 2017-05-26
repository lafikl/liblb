package liblb

import (
	"math/rand"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type host struct {
	name string
	load uint64
}

// P2C will distribute the traffic by choosing two random hosts
// and then pick the least loaded of the two.
//
// It guarantees that the load variance between any two servers
// will never exceed log(log(n)) where n is the number of hosts
type P2C struct {
	hosts   []*host
	rndm    *rand.Rand
	loadMap map[string]*host

	enableMetrics bool
	servedReqs    *prometheus.CounterVec

	sync.Mutex
}

// New returns a new instance of RandomTwoBalancer
func NewP2C() *P2C {
	return &P2C{
		hosts:   []*host{},
		loadMap: map[string]*host{},
		rndm:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

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

	p.enableMetrics = true
	return nil
}

func (p *P2C) AddHost(hostName string, load uint64) {
	p.Lock()
	defer p.Unlock()

	h := &host{name: hostName, load: load}
	p.hosts = append(p.hosts, h)
	p.loadMap[hostName] = h
}

func (p *P2C) RemoveHost(host ...string) {
	p.Lock()
	defer p.Unlock()

	for _, h := range host {
		_, ok := p.loadMap[h]
		if !ok {
			continue
		}

		delete(p.loadMap, h)

		for i, v := range p.hosts {
			if v.name == h {
				p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			}
		}

	}
}

// Balance picks two servers randomly then returns
// the least loaded one between the two
func (p *P2C) Balance() string {
	p.Lock()
	defer p.Unlock()

	// chosen host
	var host string

	n1 := p.hosts[p.rndm.Intn(len(p.hosts))].name
	n2 := p.hosts[p.rndm.Intn(len(p.hosts))].name

	host = n2

	if p.loadMap[n1].load <= p.loadMap[n1].load {
		host = n1
	}

	if p.enableMetrics {
		p.servedReqs.WithLabelValues(host).Inc()
	}

	p.loadMap[host].load++
	return host
}

// UpdateLoad updates the load of a host
func (p *P2C) UpdateLoad(host string, load uint64) error {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return ErrNoHost
	}
	h.load = load
	return nil
}

func (p *P2C) GetLoad(host string) (load uint64, err error) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return 0, ErrNoHost
	}
	return h.load, nil
}
