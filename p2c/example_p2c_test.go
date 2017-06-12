package p2c_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/lafikl/liblb/p2c"
)

func Example(t *testing.T) {
	hosts := []string{"127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007"}

	// Power of Two choices example
	lb := p2c.New(hosts...)
	for i := 0; i < 10; i++ {
		// uses random power of two choices, because the key length == 0
		host, err := lb.Balance("")
		if err != nil {
			log.Fatal(err)
		}
		// load should be around 33% per host
		fmt.Printf("Send request #%d to host %s\n", i, host)
		// when the work assign to the host is done
		lb.Done(host)
	}

	// Partial Key Grouping example
	pp2c := p2c.New(hosts...)
	for i := 0; i < 10; i++ {
		// uses PKG because the key length is > 0
		host, err := pp2c.Balance("hello world")
		if err != nil {
			log.Fatal(err)
		}

		// traffic should be split between two nodes only
		fmt.Printf("Send request #%d to host %s\n", i, host)
		// when the work assign to the host is done
		pp2c.Done(host)

	}
}
