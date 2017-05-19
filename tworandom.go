package liblb

import (
	"math/rand"
	"time"
)

type server struct {
	name string
	load uint64
}

// RandomTwoBalancer implements alogrithm described in the following paper
// https://www.eecs.harvard.edu/~michaelm/postscripts/mythesis.pdf
type RandomTwoBalancer struct {
	servers []server
	rndm    *rand.Rand
	loadMap map[string]uint64
}

// New returns a new instance of RandomTwoBalancer
func New(servers []string) *RandomTwoBalancer {
	s := []server{}
	for _, item := range servers {
		s = append(s, server{name: item, load: 0})
	}
	return &RandomTwoBalancer{
		servers: s,
		rndm:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Balance picks two servers randomly then returns
// the least loaded one between the two
func (rt *RandomTwoBalancer) Balance() string {
	n1 := rt.servers[rt.rndm.Intn(len(rt.servers))].name
	n2 := rt.servers[rt.rndm.Intn(len(rt.servers))].name

	if rt.loadMap[n1] <= rt.loadMap[n1] {
		return n1
	}
	return n2
}

func (rt *RandomTwoBalancer) AddLoad(server string, load uint64) {
	rt.loadMap[server] = load
}

func (rt *RandomTwoBalancer) GetLoad() map[string]uint64 {
	return rt.loadMap
}
