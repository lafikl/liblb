package liblb

import "stathat.com/c/consistent"

type Consistent struct {
	ch *consistent.Consistent
}

func NewConsistent() *Consistent {
	return &Consistent{ch: consistent.New()}
}

func (h *Consistent) AddHost(host string) {
	h.ch.Add(host)
}

func (h *Consistent) RemoveHost(host string) {
	h.ch.Remove(host)
}

func (h *Consistent) Balance(key string) (host string, err error) {
	return h.ch.Get(key)
}
