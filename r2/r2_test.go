package r2

import (
	"log"
	"testing"
)

func TestNewR2(t *testing.T) {
	hosts := []string{"127.0.0.1", "94.0.0.1", "88.0.0.1"}
	reqPerHost := 100

	lb := New(hosts...)
	loads := map[string]uint64{}

	for i := 0; i < reqPerHost*len(hosts); i++ {
		host, _ := lb.Balance()

		l, _ := loads[host]
		loads[host] = l + 1
	}
	for h, load := range loads {
		if load > uint64(reqPerHost) {
			t.Fatalf("host(%s) got overloaded %d>%d\n", h, load, reqPerHost)
		}
	}
	log.Println(loads)
}

func TestWeightedR2(t *testing.T) {
	hosts := []string{"127.0.0.1", "94.0.0.1", "88.0.0.1"}
	reqPerHost := 100

	lb := New()

	// in reverse order just to make sure
	// that insetion order of hosts doesn't affect anything
	for i := len(hosts); i > 0; i-- {
		lb.AddWeight(hosts[i-1], i)
	}

	loads := map[string]uint64{}

	for i := 0; i < reqPerHost*len(hosts); i++ {
		host, _ := lb.Balance()

		l, _ := loads[host]
		loads[host] = l + 1
	}

	for i, host := range hosts {
		if loads[host] > uint64(reqPerHost*(i+1)) {
			t.Fatalf("host(%s) got overloaded %d>%d\n", host, loads[host], reqPerHost*i)
		}
	}
}
