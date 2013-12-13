package utils

import (
	"path/filepath"
	"reflect"
	"testing"
)

var noExistedFile = "/tmp/not_existed_file"

func TestSelfPath(t *testing.T) {
	path := SelfPath()
	if path == "" {
		t.Error("path cannot be empty")
	}
	t.Logf("SelfPath: %s", path)
}

func TestSelfDir(t *testing.T) {
	dir := SelfDir()
	t.Logf("SelfDir: %s", dir)
}

func TestFileExists(t *testing.T) {
	if !FileExists("./file.go") {
		t.Errorf("./file.go should exists, but it didn't")
	}

	if FileExists(noExistedFile) {
		t.Errorf("Wierd, how could this file exists: %s", noExistedFile)
	}
}

func TestSearchFile(t *testing.T) {
	path, err := SearchFile(filepath.Base(SelfPath()), SelfDir())
	if err != nil {
		t.Error(err)
	}
	t.Log(path)

	path, err = SearchFile(noExistedFile, ".")
	if err == nil {
		t.Errorf("err shouldnot be nil, got path: %s", SelfDir())
	}
}

func TestGrepFile(t *testing.T) {
	_, err := GrepFile("", noExistedFile)
	if err == nil {
		t.Error("expect file-not-existed error, but got nothing")
	}

	path := filepath.Join(".", "testdata", "grepe.test")
	lines, err := GrepFile(`^\s*[^#]+`, path)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(lines, []string{"hello", "world"}) {
		t.Errorf("expect [hello world], but receive %v", lines)
	}
}
