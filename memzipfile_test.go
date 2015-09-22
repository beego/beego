package beego

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const licenseFile = "./LICENSE"

func TestOpenMemZipFile_1(t *testing.T) {
	mf, err := openMemZipFile(licenseFile, "")
	if err != nil {
		t.Fail()
	}
	file, _ := os.Open(licenseFile)
	content, _ := ioutil.ReadAll(file)
	assetMenFileAndContent(mf, content, t)
}

func assetMenFileAndContent(mf *memFile, content []byte, t *testing.T) {
	if mf.fi.contentSize != int64(len(content)) {
		t.Log("content size not same")
		t.Fail()
	}
	for i, v := range content {
		if v != mf.fi.content[i] {
			t.Log("content not same")
			t.Fail()
		}
	}
	if len(menFileInfoMap) == 0 {
		t.Log("men map is empty")
		t.Fail()
	}
}
func TestOpenMemZipFile_2(t *testing.T) {
	mf, err := openMemZipFile(licenseFile, "gzip")
	if err != nil {
		t.Fail()
	}
	file, _ := os.Open(licenseFile)
	var zipBuf bytes.Buffer
	fileWriter, _ := gzip.NewWriterLevel(&zipBuf, gzip.BestCompression)
	io.Copy(fileWriter, file)
	fileWriter.Close()
	content, _ := ioutil.ReadAll(&zipBuf)
	assetMenFileAndContent(mf, content, t)
}
func TestOpenMemZipFile_3(t *testing.T) {
	mf, err := openMemZipFile(licenseFile, "deflate")
	if err != nil {
		t.Fail()
	}
	file, _ := os.Open(licenseFile)
	var zipBuf bytes.Buffer
	fileWriter, _ := flate.NewWriter(&zipBuf, flate.BestCompression)
	io.Copy(fileWriter, file)
	fileWriter.Close()
	content, _ := ioutil.ReadAll(&zipBuf)
	assetMenFileAndContent(mf, content, t)
}
