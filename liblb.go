package liblb

import "errors"

// Balancer should be used when you need to
// accept a certain load balancer as an argument
//
// candidate implementations:
//    - roundrobin
//    - hashing/sticky
//
// cache friendly algorithms
//    - Hashing
//    - Consistent Hashing
//    - Consistent Hashing with Bounded Loads
//
// Least Loaded algorithms:
//     - TwoRandomChoices
//     - EWMA
//
//     Load can be anything, for example:
//         - Connections
//         - Latency
//         - Traffic
type Balancer interface {
	Balance() string
}

var ErrNoHost = errors.New("host not found")
