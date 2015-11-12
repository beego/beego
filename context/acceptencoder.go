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

package context

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"io/ioutil"
)

// WriteFile reads from file and writes to writer by the specific encoding(gzip/deflate)

func WriteFile(encoding string, writer io.Writer, file *os.File) (bool, string, error) {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return false, "", err
	}
	return writeLevel(encoding, writer, content, flate.BestCompression)
}

// WriteBody reads  writes content to writer by the specific encoding(gzip/deflate)

func WriteBody(encoding string, writer io.Writer, content []byte) (bool, string, error) {
	return writeLevel(encoding, writer, content, flate.BestSpeed)
}

// writeLevel reads from reader,writes to writer by specific encoding and compress level
// the compress level is defined by deflate package

func writeLevel(encoding string, writer io.Writer, content []byte, level int) (bool, string, error) {
	var outputWriter io.Writer
	var err error
	if cf, ok := encoderMap[encoding]; ok {
		outputWriter, err = cf.encode(writer, level)
	} else {
		encoding = ""
		outputWriter, err = noneCompress(writer, level)
	}
	if err != nil {
		return false, "", err
	}
	outputWriter.Write(content)
	switch outputWriter.(type) {
	case io.WriteCloser:
		outputWriter.(io.WriteCloser).Close()
	}
	return encoding != "", encoding, nil
}

// ParseEncoding will extract the right encoding for response
// the Accept-Encoding's sec is here:
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.3

func ParseEncoding(r *http.Request) string {
	if r == nil {
		return ""
	}
	return parseEncoding(r)
}

type q struct {
	name  string
	value float64
}

func noneCompress(wr io.Writer, level int) (io.Writer, error) {
	return wr, nil
}
func gzipCompress(wr io.Writer, level int) (io.Writer, error) {
	return gzip.NewWriterLevel(wr, level)
}
func deflateCompress(wr io.Writer, level int) (io.Writer, error) {
	return flate.NewWriter(wr, level)
}

type acceptEncoder struct {
	name   string
	encode func(io.Writer, int) (io.Writer, error)
}

var (
	encoderMap = map[string]acceptEncoder{ // all the other compress methods will ignore
		"gzip":     acceptEncoder{"gzip", gzipCompress},
		"deflate":  acceptEncoder{"deflate", deflateCompress},
		"*":        acceptEncoder{"gzip", gzipCompress}, // * means any compress will accept,we prefer gzip
		"identity": acceptEncoder{"", noneCompress},     // identity means none-compress
	}
)

func parseEncoding(r *http.Request) string {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	if acceptEncoding == "" {
		return ""
	}
	var lastQ q
	for _, v := range strings.Split(acceptEncoding, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		vs := strings.Split(v, ";")
		if len(vs) == 1 {
			lastQ = q{vs[0], 1}
			break
		}
		if len(vs) == 2 {
			f, _ := strconv.ParseFloat(strings.Replace(vs[1], "q=", "", -1), 64)
			if f == 0 {
				continue
			}
			if f > lastQ.value {
				lastQ = q{vs[0], f}
			}
		}
	}
	if cf, ok := encoderMap[lastQ.name]; ok {
		return cf.name
	} else {
		return ""
	}
}
