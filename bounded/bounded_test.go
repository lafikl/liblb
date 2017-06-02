package bounded

import (
	"fmt"
	"testing"
)

func TestNewConsistentBounded(t *testing.T) {
	lb := New("127.0.0.1", "192.0.0.1", "88.0.0.1", "10.0.0.1")

	for i := 0; i < 1000; i++ {
		// _, err := lb.Balance(fmt.Sprintf("hello world %d", i))
		_, err := lb.Balance("hello world")
		if err != nil {
			t.Fatal(err)
		}
	}

	loads := lb.Loads()
	for k, load := range loads {
		if load > lb.MaxLoad(k) {
			t.Fatal(fmt.Sprintf("%s load(%d) > avgLoad(%d)", k,
				load, lb.AvgLoad()))
		}
	}
	for k, load := range loads {
		fmt.Printf("%s load(%d,%d) > avgLoad(%d)\n", k, lb.MaxLoad(k),
			load, lb.AvgLoad())
	}
}
