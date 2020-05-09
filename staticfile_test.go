package beego

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var currentWorkDir, _ = os.Getwd()
var licenseFile = filepath.Join(currentWorkDir, "LICENSE")

func testOpenFile(encoding string, content []byte, t *testing.T) {
	fi, _ := os.Stat(licenseFile)
	b, n, sch, reader, err := openFile(licenseFile, fi, encoding)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log("open static file encoding "+n, b)

	assetOpenFileAndContent(sch, reader, content, t)
}
func TestOpenStaticFile_1(t *testing.T) {
	file, _ := os.Open(licenseFile)
	content, _ := ioutil.ReadAll(file)
	testOpenFile("", content, t)
}

func TestOpenStaticFileGzip_1(t *testing.T) {
	file, _ := os.Open(licenseFile)
	var zipBuf bytes.Buffer
	fileWriter, _ := gzip.NewWriterLevel(&zipBuf, gzip.BestCompression)
	io.Copy(fileWriter, file)
	fileWriter.Close()
	content, _ := ioutil.ReadAll(&zipBuf)

	testOpenFile("gzip", content, t)
}
func TestOpenStaticFileDeflate_1(t *testing.T) {
	file, _ := os.Open(licenseFile)
	var zipBuf bytes.Buffer
	fileWriter, _ := zlib.NewWriterLevel(&zipBuf, zlib.BestCompression)
	io.Copy(fileWriter, file)
	fileWriter.Close()
	content, _ := ioutil.ReadAll(&zipBuf)

	testOpenFile("deflate", content, t)
}

func TestStaticCacheWork(t *testing.T) {
	encodings := []string{"", "gzip", "deflate"}

	fi, _ := os.Stat(licenseFile)
	for _, encoding := range encodings {
		_, _, first, _, err := openFile(licenseFile, fi, encoding)
		if err != nil {
			t.Error(err)
			continue
		}

		_, _, second, _, err := openFile(licenseFile, fi, encoding)
		if err != nil {
			t.Error(err)
			continue
		}

		address1 := fmt.Sprintf("%p", first)
		address2 := fmt.Sprintf("%p", second)
		if address1 != address2 {
			t.Errorf("encoding '%v' can not hit cache", encoding)
		}
	}
}

func assetOpenFileAndContent(sch *serveContentHolder, reader *serveContentReader, content []byte, t *testing.T) {
	t.Log(sch.size, len(content))
	if sch.size != int64(len(content)) {
		t.Log("static content file size not same")
		t.Fail()
	}
	bs, _ := ioutil.ReadAll(reader)
	for i, v := range content {
		if v != bs[i] {
			t.Log("content not same")
			t.Fail()
		}
	}
	if staticFileLruCache.Len() == 0 {
		t.Log("men map is empty")
		t.Fail()
	}
}
