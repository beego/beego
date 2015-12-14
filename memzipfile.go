// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beego

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var gmfim map[string]*memFileInfo = make(map[string]*memFileInfo)
var lock sync.RWMutex

// OpenMemZipFile returns MemFile object with a compressed static file.
// it's used for serve static file if gzip enable.
func openMemZipFile(path string, zip string) (*memFile, error) {
	osfile, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer osfile.Close()

	osfileinfo, e := osfile.Stat()
	if e != nil {
		return nil, e
	}

	modtime := osfileinfo.ModTime()
	fileSize := osfileinfo.Size()
	lock.RLock()
	cfi, ok := gmfim[zip+":"+path]
	lock.RUnlock()
	if !(ok && cfi.ModTime() == modtime && cfi.fileSize == fileSize) {
		var content []byte
		if zip == "gzip" {
			var zipbuf bytes.Buffer
			gzipwriter, e := gzip.NewWriterLevel(&zipbuf, gzip.BestCompression)
			if e != nil {
				return nil, e
			}
			_, e = io.Copy(gzipwriter, osfile)
			gzipwriter.Close()
			if e != nil {
				return nil, e
			}
			content, e = ioutil.ReadAll(&zipbuf)
			if e != nil {
				return nil, e
			}
		} else if zip == "deflate" {
			var zipbuf bytes.Buffer
			deflatewriter, e := flate.NewWriter(&zipbuf, flate.BestCompression)
			if e != nil {
				return nil, e
			}
			_, e = io.Copy(deflatewriter, osfile)
			deflatewriter.Close()
			if e != nil {
				return nil, e
			}
			content, e = ioutil.ReadAll(&zipbuf)
			if e != nil {
				return nil, e
			}
		} else {
			content, e = ioutil.ReadAll(osfile)
			if e != nil {
				return nil, e
			}
		}

		cfi = &memFileInfo{osfileinfo, modtime, content, int64(len(content)), fileSize}
		lock.Lock()
		defer lock.Unlock()
		gmfim[zip+":"+path] = cfi
	}
	return &memFile{fi: cfi, offset: 0}, nil
}

// MemFileInfo contains a compressed file bytes and file information.
// it implements os.FileInfo interface.
type memFileInfo struct {
	os.FileInfo
	modTime     time.Time
	content     []byte
	contentSize int64
	fileSize    int64
}

// Name returns the compressed filename.
func (fi *memFileInfo) Name() string {
	return fi.Name()
}

// Size returns the raw file content size, not compressed size.
func (fi *memFileInfo) Size() int64 {
	return fi.contentSize
}

// Mode returns file mode.
func (fi *memFileInfo) Mode() os.FileMode {
	return fi.Mode()
}

// ModTime returns the last modified time of raw file.
func (fi *memFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir returns the compressing file is a directory or not.
func (fi *memFileInfo) IsDir() bool {
	return fi.IsDir()
}

// return nil. implement the os.FileInfo interface method.
func (fi *memFileInfo) Sys() interface{} {
	return nil
}

// MemFile contains MemFileInfo and bytes offset when reading.
// it implements io.Reader,io.ReadCloser and io.Seeker.
type memFile struct {
	fi     *memFileInfo
	offset int64
}

// Close memfile.
func (f *memFile) Close() error {
	return nil
}

// Get os.FileInfo of memfile.
func (f *memFile) Stat() (os.FileInfo, error) {
	return f.fi, nil
}

// read os.FileInfo of files in directory of memfile.
// it returns empty slice.
func (f *memFile) Readdir(count int) ([]os.FileInfo, error) {
	infos := []os.FileInfo{}

	return infos, nil
}

// Read bytes from the compressed file bytes.
func (f *memFile) Read(p []byte) (n int, err error) {
	if len(f.fi.content)-int(f.offset) >= len(p) {
		n = len(p)
	} else {
		n = len(f.fi.content) - int(f.offset)
		err = io.EOF
	}
	copy(p, f.fi.content[f.offset:f.offset+int64(n)])
	f.offset += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

// Read bytes from the compressed file bytes by seeker.
func (f *memFile) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	default:
		return 0, errWhence
	case os.SEEK_SET:
	case os.SEEK_CUR:
		offset += f.offset
	case os.SEEK_END:
		offset += int64(len(f.fi.content))
	}
	if offset < 0 || int(offset) > len(f.fi.content) {
		return 0, errOffset
	}
	f.offset = offset
	return f.offset, nil
}

// GetAcceptEncodingZip returns accept encoding format in http header.
// zip is first, then deflate if both accepted.
// If no accepted, return empty string.
func getAcceptEncodingZip(r *http.Request) string {
	ss := r.Header.Get("Accept-Encoding")
	ss = strings.ToLower(ss)
	if strings.Contains(ss, "gzip") {
		return "gzip"
	} else if strings.Contains(ss, "deflate") {
		return "deflate"
	} else {
		return ""
	}
}
