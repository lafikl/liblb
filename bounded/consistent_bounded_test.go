package bounded

import (
	"fmt"
	"testing"
)

func TestNewConsistentBounded(t *testing.T) {
	lb := New("127.0.0.1", "192.0.0.1", "88.0.0.1", "10.0.0.1")

	for i := 0; i < 10*100; i++ {
		_, err := lb.Balance(fmt.Sprintf("hello world %d", i))
		if err != nil {
			t.Fatal(err)
		}
	}

	loads := lb.Loads()
	weights := lb.Weights()
	for k, load := range loads {
		if load > lb.AvgLoad()*weights[k] {
			t.Fatal(fmt.Sprintf("%s load(%d) > avgLoad(%d)", k,
				load, lb.AvgLoad()*weights[k]))
		}
	}
	for k, load := range loads {
		fmt.Printf("%s load(%d,%d) > avgLoad(%d)\n", k, lb.MaxLoad(k),
			load, lb.AvgLoad()*weights[k])
	}
}
