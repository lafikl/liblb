package liblb

import (
	"math/rand"
	"time"
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
	hosts   []host
	rndm    *rand.Rand
	loadMap map[string]*host
}

// New returns a new instance of RandomTwoBalancer
func NewP2C(hosts []string) *P2C {
	h := []host{}
	for _, item := range hosts {
		h = append(h, host{name: item, load: 0})
	}
	return &P2C{
		hosts: h,
		rndm:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *P2C) AddHost(hostName string, load uint64) {
	p.hosts = append(p.hosts, host{name: hostName, load: load})
}

func (p *P2C) RemoveHost(host ...string) {
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
	n1 := p.hosts[p.rndm.Intn(len(p.hosts))].name
	n2 := p.hosts[p.rndm.Intn(len(p.hosts))].name

	if p.loadMap[n1].load <= p.loadMap[n1].load {
		p.loadMap[n1].load++
		return n1
	}
	p.loadMap[n2].load++
	return n2
}

// UpdateLoad updates the load of a host
func (p *P2C) UpdateLoad(host string, load uint64) error {
	h, ok := p.loadMap[host]
	if !ok {
		return ErrNoHost
	}
	h.load = load
	return nil
}

func (p *P2C) GetLoad(host string) (load uint64, err error) {
	h, ok := p.loadMap[host]
	if !ok {
		return 0, ErrNoHost
	}
	return h.load, nil
}
