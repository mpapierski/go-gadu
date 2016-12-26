package gadu

import "testing"

func TestVersion(t *testing.T) {
	version := Version()
	if version != "1.12.1" {
		t.Errorf("Invalid version: %s", version)
	}
}
