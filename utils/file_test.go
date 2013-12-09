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
	if !FileExists("/bin/echo") {
		t.Errorf("/bin/echo should exists, but it didn't")
	}

	if FileExists(noExistedFile) {
		t.Errorf("Wierd, how could this file exists: %s", noExistedFile)
	}
}

func TestLookFile(t *testing.T) {
	path, err := LookFile(filepath.Base(SelfPath()), SelfDir())
	if err != nil {
		t.Error(err)
	}
	t.Log(path)

	path, err = LookFile(noExistedFile, ".")
	if err == nil {
		t.Errorf("err shouldnot be nil, got path: %s", SelfDir())
	}
}

func TestGrepE(t *testing.T) {
	_, err := GrepE("", noExistedFile)
	if err == nil {
		t.Error("expect file-not-existed error, but got nothing")
	}

	path := filepath.Join(".", "testdata", "grepe.test")
	lines, err := GrepE(`^\s*[^#]+`, path)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(lines, []string{"hello", "world"}) {
		t.Errorf("expect [hello world], but receive %v", lines)
	}
}
