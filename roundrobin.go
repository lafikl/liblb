package liblb

// RoundRobin implements round-robin alogrithm
type RoundRobin struct {
	i       int
	servers []string
}

func NewRoundRobin(servers []string) *RoundRobin {
	return &RoundRobin{i: 0, servers: servers}
}

func (rb *RoundRobin) Balance() string {
	server := rb.servers[rb.i]
	rb.i++
	if rb.i >= len(rb.servers) {
		rb.i = 0
	}
	return server
}
