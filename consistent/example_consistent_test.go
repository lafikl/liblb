package consistent_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/lafikl/liblb/consistent"
)

func Example(t *testing.T) {
	hosts := []string{"127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007"}

	lb := consistent.New(hosts...)
	for i := 0; i < 10; i++ {
		host, err := lb.Balance(fmt.Sprintf("hello world %d", i))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Send request #%d to host %s\n", i, host)
	}
}
