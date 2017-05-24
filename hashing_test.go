package liblb

import "testing"

func TestNewHash(t *testing.T) {
	lb := NewHashBalancer()
	lb.AddHost("127.0.0.1")
}
