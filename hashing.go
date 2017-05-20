package liblb

import "stathat.com/c/consistent"

type HashBalancer struct {
	ch *consistent.Consistent
}

func NewHashBalancer() *HashBalancer {
	return &HashBalancer{ch: consistent.New()}
}

func (h *HashBalancer) AddHost(host string) {
	h.ch.Add(host)
}

func (h *HashBalancer) RemoveHost(host string) {
	h.ch.Remove(host)
}

func (h *HashBalancer) Balance(key string) (host string, err error) {
	return h.ch.Get(key)
}
