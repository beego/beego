package utils

import (
	"os"
	"path/filepath"
)

func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}
