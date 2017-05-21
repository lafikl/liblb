package liblb

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestNewP2C(t *testing.T) {
	var hosts = []string{
		"127.0.0.1",
		"225.0.0.1",
		"10.0.0.1",
		"28.0.0.1",
		"88.0.0.1",
	}
	upperVariance := 1 - math.Log(math.Log(float64(len(hosts))))

	lb := NewP2C()

	for _, host := range hosts {
		lb.AddHost(host, 0)
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
		floatVal := float64(val)
		fmt.Printf("%s=%d\n", host, val)

		// check for load variance

		for _, otherHost := range hosts {
			oval, err := lb.GetLoad(otherHost)
			if err != nil {
				t.Fatal(err)
			}
			floatOval := float64(oval)

			variance := floatVal / (floatVal + floatOval)
			if variance > goalVariance {
				t.Fatalf("variance between (%s, %s) is %.2f > %.2f\n",
					host, otherHost, variance, upperVariance)
			} else {
				fmt.Printf("variance between (%s, %s) is %.2f and upper is %.2f\n",
					host, otherHost, variance, upperVariance)
			}
		}
	}

}
