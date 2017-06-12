package p2c

import (
	"fmt"
	"log"
	"math"
	"testing"
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

	lb := New()

	for _, host := range hosts {
		lb.Add(host)
	}

	for i := 0; i < 200; i++ {
		lb.Balance("")
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
			if variance > upperVariance {
				t.Fatalf("variance between (%s, %s) is %.2f > %.2f\n",
					host, otherHost, variance, upperVariance)
			} else {
				fmt.Printf("variance between (%s, %s) is %.2f and upper is %.2f\n",
					host, otherHost, variance, upperVariance)
			}
		}
	}
}

func TestNewHP2C(t *testing.T) {
	var hosts = []string{
		"127.0.0.1",
		"225.0.0.1",
		"10.0.0.1",
		"28.0.0.1",
		"88.0.0.1",
	}

	lb := New()

	for _, host := range hosts {
		lb.Add(host)
	}

	for i := 0; i < 200; i++ {
		for j := 0; j < 100; j++ {
			lb.Balance("hello, world!")
		}
	}

	for _, host := range hosts {
		val, err := lb.GetLoad(host)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s=%d\n", host, val)

	}
}

func TestLongestRun(t *testing.T) {
	var hosts = []string{
		"127.0.0.1",
		"225.0.0.1",
		"10.0.0.1",
		"28.0.0.1",
		"88.0.0.1",
	}
	lb := New(hosts...)

	longestRun := map[string]int{}

	currentHost := ""
	currentCount := 0
	for i := 0; i < 1000; i++ {
		// host, err := lb.Balance(fmt.Sprintf("hello, world !", i))
		host, err := lb.Balance("hello, world!")
		if err != nil {
			t.Fatal(err)
		}
		if currentHost != host {
			lrun, _ := longestRun[currentHost]
			if lrun < currentCount {
				longestRun[currentHost] = currentCount
			}
			currentHost = host
			currentCount = 1
			continue
		}
		currentCount++
	}

	log.Println(longestRun)

	for _, host := range hosts {
		val, err := lb.GetLoad(host)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s=%d\n", host, val)
	}

}
