package web

import (
	"net/http"
	"os"
	"path/filepath"
)

type FileSystem struct{}

func (d FileSystem) Open(name string) (http.File, error) {
	return os.Open(name)
}

// Walk walks the file tree rooted at root in filesystem, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	f, err := fs.Open(root)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		err = walkFn(root, nil, err)
	} else {
		err = walk(fs, root, info, walkFn)
	}
	if err == filepath.SkipDir {
		return nil
	}
	return err
}

// walk recursively descends path, calling walkFn.
func walk(fs http.FileSystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	var err error
	if !info.IsDir() {
		return walkFn(path, info, nil)
	}

	dir, err := fs.Open(path)
	if err != nil {
		if err1 := walkFn(path, info, err); err1 != nil {
			return err1
		}
		return err
	}
	defer dir.Close()
	dirs, err := dir.Readdir(-1)
	err1 := walkFn(path, info, err)
	// If err != nil, walk can't walk into this directory.
	// err1 != nil means walkFn want walk to skip this directory or stop walking.
	// Therefore, if one of err and err1 isn't nil, walk will return.
	if err != nil || err1 != nil {
		// The caller's behavior is controlled by the return value, which is decided
		// by walkFn. walkFn may ignore err and return nil.
		// If walkFn returns SkipDir, it will be handled by the caller.
		// So walk should return whatever walkFn returns.
		return err1
	}

	for _, fileInfo := range dirs {
		filename := filepath.Join(path, fileInfo.Name())
		if err = walk(fs, filename, fileInfo, walkFn); err != nil {
			if !fileInfo.IsDir() || err != filepath.SkipDir {
				return err
			}
		}
	}
	return nil
}
