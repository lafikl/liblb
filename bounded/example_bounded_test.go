package bounded_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/lafikl/liblb/bounded"
)

func Example(t *testing.T) {
	hosts := []string{"127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007"}

	lb := bounded.New(hosts...)
	// Host load will never exceed (number_of_requests/len(hosts)) by more than 25%
	// in this case:
	// any host load would be at most:
	// ceil((10/3) * 1.25)
	for i := 0; i < 10; i++ {
		host, err := lb.Balance("hello world")
		if err != nil {
			log.Fatal(err)
		}
		// do work for "host"
		fmt.Printf("Send request #%d to host %s\n", i, host)
		// when the work assign to the host is done
		lb.Done(host)
	}

}
