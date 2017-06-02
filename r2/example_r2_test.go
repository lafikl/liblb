package r2_test

import (
	"fmt"
	"testing"

	"github.com/lafikl/liblb/r2"
)

func Example(t *testing.T) {
	lb := r2.New()
	lb.AddHost("127.0.0.1:8009")
	fmt.Println("Hello")
}
