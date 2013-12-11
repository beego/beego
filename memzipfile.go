//借鉴了 https://github.com/ranveerkunal/memfs
package beego

/**
* 思路如下:
* fs.go:
type FileInfo interface {
    Name() string       // base name of the file
    Size() int64        // length in bytes for regular files; system-dependent for others
    Mode() FileMode     // file mode bits
    ModTime() time.Time // modification time
    IsDir() bool        // abbreviation for Mode().IsDir()
    Sys() interface{}   // underlying data source (can return nil)
}

// A FileSystem implements access to a collection of named files.
// The elements in a file path are separated by slash ('/', U+002F)
// characters, regardless of host operating system convention.
type FileSystem interface {
	Open(name string) (File, error)
}

// A File is returned by a FileSystem's Open method and can be
// served by the FileServer implementation.
type File interface {
	Close() error
	Stat() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
	Read([]byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
}

func ServeFile(w ResponseWriter, r *Request, name string) {
	dir, file := filepath.Split(name)
	serveFile(w, r, Dir(dir), file, false)
}

//不考虑, dir和redirect的情况
func serveFile(w ResponseWriter, r *Request, fs FileSystem, name string) {

	f, err := fs.Open(name)
	if err != nil {
		// TODO expose actual error?
		NotFound(w, r)
		return
	}
	defer f.Close()

	d, err1 := f.Stat()
	if err1 != nil {
		// TODO expose actual error?
		NotFound(w, r)
		return
	}

	// serverContent will check modification time
	serveContent(w, r, d.Name(), d.ModTime(), d.Size(), f)
}

*/

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	//"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var gmfim map[string]*MemFileInfo = make(map[string]*MemFileInfo)

//TODO: 加锁保证数据完整性
func OpenMemZipFile(path string, zip string) (*MemFile, error) {
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

	cfi, ok := gmfim[zip+":"+path]
	if ok && cfi.ModTime() == modtime && cfi.fileSize == fileSize {
		//fmt.Printf("read %s file %s from cache\n", zip, path)
	} else {
		//fmt.Printf("NOT read %s file %s from cache\n", zip, path)
		var content []byte
		if zip == "gzip" {
			//将文件内容压缩到zipbuf中
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
			//读zipbuf到content
			content, e = ioutil.ReadAll(&zipbuf)
			if e != nil {
				return nil, e
			}
		} else if zip == "deflate" {
			//将文件内容压缩到zipbuf中
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
			//将zipbuf读入到content
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

		cfi = &MemFileInfo{osfileinfo, modtime, content, int64(len(content)), fileSize}
		gmfim[zip+":"+path] = cfi
		//fmt.Printf("%s file %s to %d, cache it\n", zip, path, len(content))
	}
	return &MemFile{fi: cfi, offset: 0}, nil
}

type MemFileInfo struct {
	os.FileInfo
	modTime     time.Time
	content     []byte
	contentSize int64
	fileSize    int64
}

func (fi *MemFileInfo) Name() string {
	return fi.Name()
}

func (fi *MemFileInfo) Size() int64 {
	return fi.contentSize
}

func (fi *MemFileInfo) Mode() os.FileMode {
	return fi.Mode()
}

func (fi *MemFileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *MemFileInfo) IsDir() bool {
	return fi.IsDir()
}

func (fi *MemFileInfo) Sys() interface{} {
	return nil
}

type MemFile struct {
	fi     *MemFileInfo
	offset int64
}

func (f *MemFile) Close() error {
	return nil
}

func (f *MemFile) Stat() (os.FileInfo, error) {
	return f.fi, nil
}

func (f *MemFile) Readdir(count int) ([]os.FileInfo, error) {
	infos := []os.FileInfo{}

	return infos, nil
}

func (f *MemFile) Read(p []byte) (n int, err error) {
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

func (f *MemFile) Seek(offset int64, whence int) (ret int64, err error) {
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

//返回: gzip, deflate, 优先gzip
//返回空, 表示不zip
func GetAcceptEncodingZip(r *http.Request) string {
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

func CloseZWriter(zwriter io.Writer) {
	if zwriter == nil {
		return
	}

	switch zwriter.(type) {
	case *gzip.Writer:
		zwriter.(*gzip.Writer).Close()
	case *flate.Writer:
		zwriter.(*flate.Writer).Close()
		//其他情况不close, 保持和默认(非压缩)行为一致
		/*
			case io.WriteCloser:
				zwriter.(io.WriteCloser).Close()
		*/
	}
}
