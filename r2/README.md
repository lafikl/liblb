# Round Robin Balancer
It's one of the simplest if not the simplest load balancing algorithm.
It distribute requests by walking the list of servers and assigning a request to each server in turn.
On the downside Round-Robin assumes that all servers are alike,
and that all requests take the same amount of time, which is obviously not true in practice.


### Usage Example:
```go
package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/lafikl/liblb/r2"
)

func main() {
	lb := r2.New("127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007")
	for i := 0; i < 10; i++ {
		host, err := lb.Balance()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send request #%d to host %s\n", i, host)
	}
}
```


## Weighted Round Robin
A variant of Round-Robin that assigns a weight for every host, the weight affects the number of requests that gets sent to the server.
Assume that we have two hosts `A` with weight **1** and `B` with weight **4**,
that means for every single request we send to `A` we send 4 requests to `B`.
In other words, 80% of the requests would go to `B`.
Which you can calculate by yourself applying this formula `host_weight/total_weights`.


### Usage Example:
```go
package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/lafikl/liblb/r2"
)

func main() {
    // default weight is 1
	lb := r2.New("127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007")
    // host_weight/total_weights
    // this hosts load would be 3/(3+3)=0.5
    // meaning that 50% of the requests would go to 127.0.0.1:9000
    lb.AddWeight("127.0.0.1:9000", 3)
	for i := 0; i < 10; i++ {
		host, err := lb.Balance()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send request #%d to host %s\n", i, host)
	}
}
```
