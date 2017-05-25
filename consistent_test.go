package liblb

import "testing"

func TestNewConsistent(t *testing.T) {
	lb := NewConsistent()
	lb.AddHost("127.0.0.1")
	// @TODO(kl): add real tests
}
