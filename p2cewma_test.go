package liblb

import (
	"fmt"
	"testing"
	"time"
)

func TestNewP2CEWMA(t *testing.T) {
	var hosts = []string{
		"127.0.0.1",
		"225.0.0.1",
		"10.0.0.1",
		"28.0.0.1",
		"88.0.0.1",
	}
	lb := NewP2CEWMA(1)

	for _, host := range hosts {
		lb.AddHost(host)
	}

	for i := 0; i < 200; i++ {
		for j := 0; j < 100; j++ {
			lb.Balance()
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, host := range hosts {
		val, err := lb.GetLoad(host)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s=%f\n", host, val)
	}
	lb.Close()
	lb.Stats()
}
