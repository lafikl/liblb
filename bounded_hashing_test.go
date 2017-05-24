package liblb

import (
	"fmt"
	"log"
	"testing"
)

func TestNewBoundedHash(t *testing.T) {
	lb := NewBoundedHashBalancer()
	lb.Add("127.0.0.1")
	lb.Add("192.0.0.1")
	lb.Add("88.0.0.1")
	lb.Add("10.10.0.1")
	for i := 0; i < 10*4; i++ {
		h, err := lb.Balance("hello world")
		log.Println("iter", i, h, err)
	}

	loads := lb.Loads()
	for k, load := range loads {
		if load > lb.AvgLoad() {
			t.Fatal(fmt.Sprintf("%s load(%d) > avgLoad(%d)", k, load, lb.AvgLoad()))
		}
	}
}
