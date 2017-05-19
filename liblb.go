package liblb

// Balancer should be used for
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
