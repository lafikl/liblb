package liblb

import (
	"math/rand"
	"time"

	"github.com/VividCortex/ewma"
)

type P2CEWMA struct {
	decay   float64
	hosts   []string
	loadMap map[string]ewma.MovingAverage
	rndm    *rand.Rand
}

func NewP2CEWMA(decay ...float64) *P2CEWMA {
	d := ewma.AVG_METRIC_AGE
	if len(decay) != 0 {
		d = decay[0]
	}
	e := &P2CEWMA{
		decay: d,
		rndm:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return e
}

func (p *P2CEWMA) AddHost(host ...string) {
	for _, h := range host {
		_, ok := p.loadMap[h]
		if !ok {
			p.hosts = append(p.hosts, h)
			p.loadMap[h] = ewma.NewMovingAverage(p.decay)
		}
	}
}

func (p *P2CEWMA) RemoveHost(host ...string) {
	for _, h := range host {
		_, ok := p.loadMap[h]
		if !ok {
			continue
		}

		delete(p.loadMap, h)

		for i, v := range p.hosts {
			if v == h {
				p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			}
		}

	}
}

func (p *P2CEWMA) UpdateLoad(host string, load float64) error {
	h, ok := p.loadMap[host]
	if !ok {
		return ErrNoHost
	}
	h.Add(load)
	return nil
}

func (p *P2CEWMA) GetLoad(host string) (float64, error) {
	h, ok := p.loadMap[host]
	if !ok {
		return 0, ErrNoHost
	}
	return h.Value(), nil
}

func (p *P2CEWMA) Balance() string {
	n1 := p.hosts[p.rndm.Intn(len(p.hosts))]
	n2 := p.hosts[p.rndm.Intn(len(p.hosts))]

	if p.loadMap[n1].Value() <= p.loadMap[n1].Value() {
		return n1
	}
	return n2
}
