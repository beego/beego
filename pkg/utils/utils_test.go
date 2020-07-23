package utils

import (
	"testing"
)

func TestCompareGoVersion(t *testing.T) {
	targetVersion := "go1.8"
	if compareGoVersion("go1.12.4", targetVersion) != 1 {
		t.Error("should be 1")
	}

	if compareGoVersion("go1.8.7", targetVersion) != 1 {
		t.Error("should be 1")
	}

	if compareGoVersion("go1.8", targetVersion) != 0 {
		t.Error("should be 0")
	}

	if compareGoVersion("go1.7.6", targetVersion) != -1 {
		t.Error("should be -1")
	}

	if compareGoVersion("go1.12.1rc1", targetVersion) != 1 {
		t.Error("should be 1")
	}

	if compareGoVersion("go1.8rc1", targetVersion) != 0 {
		t.Error("should be 0")
	}

	if compareGoVersion("go1.7rc1", targetVersion) != -1 {
		t.Error("should be -1")
	}
}
