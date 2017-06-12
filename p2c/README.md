# P2C Balancer
TODO


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
	lb := p2c.New("127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007")
	for i := 0; i < 10; i++ {
		host, err := lb.Balance("")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send request #%d to host %s\n", i, host)
	}
}
```

# P2C Balancer
TODO


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
	lb := p2c.New("127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007")
	for i := 0; i < 10; i++ {
		host, err := lb.Balance("")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send request #%d to host %s\n", i, host)
	}
}
```
