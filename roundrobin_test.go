package liblb

import (
	"log"
	"testing"
)

func TestNewRoundRobin(t *testing.T) {
	hosts := []string{"127.0.0.1", "94.0.0.1", "88.0.0.1"}
	reqPerHost := 100

	lb := NewRoundRobin(hosts...)
	loads := map[string]uint64{}

	for i := 0; i < reqPerHost*len(hosts); i++ {
		host := lb.Balance()

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
