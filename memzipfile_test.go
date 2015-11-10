// Copyright 2015 beego Author. All Rights Reserved.
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
