package consistent

import (
	"log"
	"testing"
)

func TestNewConsistent(t *testing.T) {
	hosts := []string{"127.0.0.1", "94.0.0.1", "88.0.0.1"}
	lb := New(hosts...)
	loads := map[string]int{}

	for i := 0; i < 100; i++ {
		host, err := lb.Balance("hello, world!")
		if err != nil {
			t.Fatal(err)
		}
		h, _ := loads[host]
		loads[host] = h + 1
	}

	// make sure that all requests got to a single host
	if len(loads) > 1 {
		t.Fatalf("load is not consistent %s\n", loads)
	}

	log.Println(loads)

}
