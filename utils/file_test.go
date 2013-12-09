package utils

import (
	"testing"
)

func TestSelfPath(t *testing.T) {
	path := SelfPath()
	if path == "" {
		t.Error("path cannot be empty")
	}
	t.Logf("SelfPath: %s", path)
}
