// Package archiver makes it super easy to create and open .zip,
// .tar.gz, and .tar.bz2 files.
package logs

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip is for Zip format
var Zip zipFormat

func init() {
	RegisterFormat("Zip", Zip)
}

type zipFormat struct{}

func (zipFormat) Match(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".zip") || isZip(filename)
}

// isZip checks the file has the Zip format signature by reading its beginning
// bytes and matching it against "PK\x03\x04"
func isZip(zipPath string) bool {
	f, err := os.Open(zipPath)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if n, err := f.Read(buf); err != nil || n < 4 {
		return false
	}

	return bytes.Equal(buf, []byte("PK\x03\x04"))
}

// Make creates a .zip file in the location zipPath containing
// the contents of files listed in filePaths. File paths
// can be those of regular files or directories. Regular
// files are stored at the 'root' of the archive, and
// directories are recursively added.
//
// Files with an extension for formats that are already
// compressed will be stored only, not compressed.
func (zipFormat) Make(zipPath string, filePaths []string) error {
	err := os.Chdir(filePaths[0])
	if err != nil {
		return err
	}

	out, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", zipPath, err)
	}
	defer out.Close()

	w := zip.NewWriter(out)
	err = zipFile(w, filePaths[0], filePaths[1])
	if err != nil {
		w.Close()
		return err
	}

	return w.Close()
}

func zipFile(w *zip.Writer, path,source string) error {
	var fpath = fmt.Sprintf("%s/%s", path, source)
	
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("%s: stat: %v", source, err)
	}
		
		
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("%s: getting header: %v", fpath, err)
	}
	header.Name = source
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("%s: making header: %v", fpath, err)
	}

	if info.IsDir() {
		return nil
	}

	if header.Mode().IsRegular() {
		file, err := os.Open(fpath)
		if err != nil {
			return fmt.Errorf("%s: opening: %v", fpath, err)
		}
		defer file.Close()

		_, err = io.CopyN(writer, file, info.Size())
		if err != nil && err != io.EOF {
			return fmt.Errorf("%s: copying contents: %v", fpath, err)
		}
	}
	
	return nil
}

// Open unzips the .zip file at source into destination.
func (zipFormat) Open(source, destination string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, zf := range r.File {
		if err := unzipFile(zf, destination); err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(zf *zip.File, destination string) error {
	if strings.HasSuffix(zf.Name, "/") {
		return mkdir(filepath.Join(destination, zf.Name))
	}

	rc, err := zf.Open()
	if err != nil {
		return fmt.Errorf("%s: open compressed file: %v", zf.Name, err)
	}
	defer rc.Close()

	return writeNewFile(filepath.Join(destination, zf.Name), rc, zf.FileInfo().Mode())
}

// compressedFormats is a (non-exhaustive) set of lowercased
// file extensions for formats that are typically already
// compressed. Compressing already-compressed files often
// results in a larger file, so when possible, we check this
// set to avoid that.
var compressedFormats = map[string]struct{}{
	".7z":   {},
	".avi":  {},
	".bz2":  {},
	".cab":  {},
	".gif":  {},
	".gz":   {},
	".jar":  {},
	".jpeg": {},
	".jpg":  {},
	".lz":   {},
	".lzma": {},
	".mov":  {},
	".mp3":  {},
	".mp4":  {},
	".mpeg": {},
	".mpg":  {},
	".png":  {},
	".rar":  {},
	".tbz2": {},
	".tgz":  {},
	".txz":  {},
	".xz":   {},
	".zip":  {},
	".zipx": {},
}
