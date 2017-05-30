package liblb

import "errors"

type Balancer interface {
	Balance() (string, error)
}

type KeyedBalancer interface {
	Balance(key string) (string, error)
}

var ErrNoHost = errors.New("host not found")
