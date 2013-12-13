package utils

import (
	"testing"
)

func TestInSlice(t *testing.T) {
	sl := []string{"A", "b"}
	if !InSlice("A", sl) {
		t.Error("should be true")
	}
	if InSlice("B", sl) {
		t.Error("should be false")
	}
}
