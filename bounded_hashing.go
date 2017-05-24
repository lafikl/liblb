package liblb

import (
	"errors"
	"fmt"
	"sync"

	"stathat.com/c/consistent"
)

var ErrAllOverloaded = errors.New("all hosts are overloaded")
var Err = errors.New("all hosts are overloaded")

type boundedHost struct {
	load   uint64
	weight int
}

type BoundedHashBalancer struct {
	ch               *consistent.Consistent
	loads            map[string]*boundedHost
	numberOfReplicas int
	totalLoad        uint64

	sync.RWMutex
}

func NewBoundedHashBalancer(numberOfReplicas ...int) *BoundedHashBalancer {
	ch := consistent.New()
	if len(numberOfReplicas) > 0 {
		ch.NumberOfReplicas = numberOfReplicas[0]
	}

	return &BoundedHashBalancer{
		ch:    consistent.New(),
		loads: map[string]*boundedHost{},
	}
}

func (b *BoundedHashBalancer) Add(host string) {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.loads[host]; ok {
		return
	}

	b.loads[host] = &boundedHost{load: 0, weight: 1}
	b.ch.Add(host)
}

func (b *BoundedHashBalancer) AddWithWeight(host string, weight int) {
	b.Lock()
	defer b.Unlock()

	if _, ok := b.loads[host]; ok {
		return
	}

	b.loads[host] = &boundedHost{load: 1, weight: weight}
	b.ch.Add(host)
}

func (b *BoundedHashBalancer) Balance(key string) (host string, err error) {
	b.Lock()
	defer b.Unlock()

	return b.get("", key, 10)
}

func (b *BoundedHashBalancer) get(firstKey, currentKey string, size int) (string, error) {
	hosts, err := b.ch.GetN(currentKey, size)
	if err != nil {
		return "", err
	}

	for _, host := range hosts {
		if host == firstKey {
			return "", ErrAllOverloaded
		}
		if b.loadOK(host) {
			b.loads[host].load++
			b.totalLoad++
			return host, nil
		}
	}
	if len(firstKey) == 0 {
		firstKey = hosts[0]
	}
	currentKey = hosts[len(hosts)-1]
	// return b.get(firstKey, currentKey, size*3/2)
	return b.get(firstKey, currentKey, size)
}

func (b *BoundedHashBalancer) Done(host string) {
	b.Lock()
	defer b.Unlock()

	bhost, ok := b.loads[host]
	if !ok {
		return
	}
	bhost.load--
	b.totalLoad--
}

func (b *BoundedHashBalancer) Loads() map[string]uint64 {
	loads := map[string]uint64{}
	for k, bhost := range b.loads {
		loads[k] = bhost.load
	}
	return loads
}

func (b *BoundedHashBalancer) loadOK(host string) bool {
	// calcs load
	if b.totalLoad == 0 {
		b.totalLoad = 1
	}
	avgLoadPerNode := b.totalLoad * 2 / 4
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	bhost, ok := b.loads[host]
	if !ok {
		panic(fmt.Sprintf("given host(%s) not in loadsMap", host))
	}

	if bhost.load < (avgLoadPerNode * uint64(bhost.weight)) {
		return true
	}

	return false
}

func (b *BoundedHashBalancer) AvgLoad() uint64 {
	avgLoadPerNode := b.totalLoad * 2 / 4
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	return avgLoadPerNode
}
