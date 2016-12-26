package gadu

import "testing"

func TestError(t *testing.T) {
	err := NewGGError(1)
	if err.Error() != "Operation not permitted" {
		t.Fatalf("Invalid error: %s", err.Error())
	}
}
