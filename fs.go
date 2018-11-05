package beego

import (
	"net/http"
	"os"
	"path/filepath"
)

type IFileSystem interface {
	http.FileSystem
	Walk(string, filepath.WalkFunc) error
}

// A File is returned by a FileSystem's Open method and can be
// served by the FileServer implementation.
//
// The methods should behave the same as those on an *os.File.
type File struct {
	*os.File
}

type FileSystem struct {
}

func (d FileSystem) Open(name string) (http.File, error) {
	return os.Open(name)
}
func (d FileSystem) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}
