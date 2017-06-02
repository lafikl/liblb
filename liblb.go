package liblb

import "errors"

type Balancer interface {
	New(hosts ...string)
	Add(host string)
	Remove(host string)
	Balance() (string, error)
}

type KeyedBalancer interface {
	New(hosts ...string)
	Add(host string)
	Remove(host string)
	Balance(key string) (string, error)
}

var ErrNoHost = errors.New("host not found")
