package liblb

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/VividCortex/ewma"
)

type peHost struct {
	name string
	load uint64
}

type P2CEWMA struct {
	decay float64

	lock      sync.RWMutex
	hosts     []*peHost
	loadMap   map[string]ewma.MovingAverage
	lastLoads map[string]uint64

	rndm      *rand.Rand
	closeChan chan struct{}
}

func NewP2CEWMA(decay ...float64) *P2CEWMA {
	d := ewma.AVG_METRIC_AGE
	if len(decay) != 0 {
		d = decay[0]
	}
	e := &P2CEWMA{
		decay:     d,
		rndm:      rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:      sync.RWMutex{},
		hosts:     []*peHost{},
		loadMap:   map[string]ewma.MovingAverage{},
		lastLoads: map[string]uint64{},
		closeChan: make(chan struct{}),
	}
	go e.trackLoad()
	return e
}

func (p *P2CEWMA) trackLoad() {
	d2 := time.Duration(int(p.decay * 1000))
	t := time.NewTicker(1 * time.Duration(d2*time.Millisecond))

loop:
	for {
		select {
		case <-t.C:
			// update our list of hosts
			p.lock.Lock()
			var currentLoad uint64
			var lastLoad uint64

			for _, host := range p.hosts {
				currentLoad = host.load
				lastLoad = p.lastLoads[host.name]

				p.loadMap[host.name].Add(float64(currentLoad - lastLoad))
				p.lastLoads[host.name] = currentLoad
			}

			p.lock.Unlock()
		case <-p.closeChan:
			t.Stop()
			break loop
		}
	}
}

func (p *P2CEWMA) AddHost(host ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, h := range host {
		_, ok := p.loadMap[h]
		if !ok {
			p.hosts = append(p.hosts, &peHost{name: h, load: 0})
			p.loadMap[h] = ewma.NewMovingAverage(p.decay)
			p.lastLoads[h] = 0
		}
	}
}

func (p *P2CEWMA) RemoveHost(host ...string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, h := range host {
		_, ok := p.loadMap[h]
		if !ok {
			continue
		}

		delete(p.loadMap, h)
		delete(p.lastLoads, h)

		for i, v := range p.hosts {
			if v.name == h {
				p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			}
		}

	}
}

func (p *P2CEWMA) UpdateLoad(host string, load float64) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return ErrNoHost
	}
	h.Add(load)
	return nil
}

func (p *P2CEWMA) GetLoad(host string) (float64, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	h, ok := p.loadMap[host]
	if !ok {
		return 0, ErrNoHost
	}
	return h.Value(), nil
}

func (p *P2CEWMA) Stats() {
	p.lock.RLock()
	defer p.lock.RUnlock()

	fmt.Println("=========Stats=============")
	for _, host := range p.hosts {
		fmt.Println(host.name, host.load)
	}
	fmt.Println("=========Stats=============")
}

func (p *P2CEWMA) Balance() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	h1 := p.hosts[p.rndm.Intn(len(p.hosts))]
	h2 := p.hosts[p.rndm.Intn(len(p.hosts))]

	if p.loadMap[h1.name].Value() <= p.loadMap[h1.name].Value() {
		h1.load++
		return h1.name
	}
	h2.load++
	return h2.name
}

func (p *P2CEWMA) Close() {
	p.closeChan <- struct{}{}
}
