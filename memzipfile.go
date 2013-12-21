package beego

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

// OpenMemZipFile returns MemFile object with a compressed static file.
// it's used for serve static file if gzip enable.
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

// MemFileInfo contains a compressed file bytes and file information.
// it implements os.FileInfo interface.
type MemFileInfo struct {
	os.FileInfo
	modTime     time.Time
	content     []byte
	contentSize int64
	fileSize    int64
}

// Name returns the compressed filename.
func (fi *MemFileInfo) Name() string {
	return fi.Name()
}

// Size returns the raw file content size, not compressed size.
func (fi *MemFileInfo) Size() int64 {
	return fi.contentSize
}

// Mode returns file mode.
func (fi *MemFileInfo) Mode() os.FileMode {
	return fi.Mode()
}

// ModTime returns the last modified time of raw file.
func (fi *MemFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir returns the compressing file is a directory or not.
func (fi *MemFileInfo) IsDir() bool {
	return fi.IsDir()
}

// return nil. implement the os.FileInfo interface method.
func (fi *MemFileInfo) Sys() interface{} {
	return nil
}

// MemFile contains MemFileInfo and bytes offset when reading.
// it implements io.Reader,io.ReadCloser and io.Seeker.
type MemFile struct {
	fi     *MemFileInfo
	offset int64
}

// Close memfile.
func (f *MemFile) Close() error {
	return nil
}

// Get os.FileInfo of memfile.
func (f *MemFile) Stat() (os.FileInfo, error) {
	return f.fi, nil
}

// read os.FileInfo of files in directory of memfile.
// it returns empty slice.
func (f *MemFile) Readdir(count int) ([]os.FileInfo, error) {
	infos := []os.FileInfo{}

	return infos, nil
}

// Read bytes from the compressed file bytes.
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

// Read bytes from the compressed file bytes by seeker.
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

// GetAcceptEncodingZip returns accept encoding format in http header.
// zip is first, then deflate if both accepted.
// If no accepted, return empty string.
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
	return ""
}

// CloseZWriter closes the io.Writer after compressing static file.
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
