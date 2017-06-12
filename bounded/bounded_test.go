package bounded

import (
	"fmt"
	"testing"
)

func TestNewConsistentBounded(t *testing.T) {
	lb := New("127.0.0.1", "192.0.0.1", "88.0.0.1", "10.0.0.1")

	for i := 0; i < 1000; i++ {
		host, err := lb.Balance(fmt.Sprintf("hello world %d", i))
		if err != nil {
			t.Fatal(err)
		}
		lb.Inc(host)
	}

	loads := lb.Loads()
	for k, load := range loads {
		if load > lb.MaxLoad() {
			t.Fatal(fmt.Sprintf("%s load(%d) > MaxLoad(%d)", k,
				load, lb.MaxLoad()))
		}
	}
	for k, load := range loads {
		fmt.Printf("%s load(%d) > MaxLoad(%d)\n", k,
			load, lb.MaxLoad())
	}
}
