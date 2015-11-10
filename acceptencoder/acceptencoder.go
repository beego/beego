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

package acceptencoder

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"bytes"
)

// WriteFile reads from file and writes to writer by the specific encoding(gzip/deflate)

func WriteFile(encoding string, writer io.Writer, file *os.File) (bool, string, error) {
	return writeLevel(encoding, writer, file, flate.BestCompression)
}

// WriteBody reads  writes content to writer by the specific encoding(gzip/deflate)

func WriteBody(encoding string, writer io.Writer, content []byte) (bool, string, error) {
	return writeLevel(encoding, writer, bytes.NewBuffer(content), flate.BestSpeed)
}

// writeLevel reads from reader,writes to writer by specific encoding and compress level
// the compress level is defined by deflate package

func writeLevel(encoding string, writer io.Writer, reader io.Reader, level int) (bool, string, error) {
	var outputWriter io.Writer
	var err error
	switch encoding {
	case "gzip":
		outputWriter, err = gzip.NewWriterLevel(writer, level)
	case "deflate":
		outputWriter, err = flate.NewWriter(writer, level)
	default:
		// all the other compress methods will ignore
		// such as the deprecated compress and chrome-only sdch
		encoding = ""
		outputWriter = writer.(io.Writer)
	}
	if err != nil {
		return false, "", err
	}
	if _, err = io.Copy(outputWriter, reader); err != nil {
		return false, "", err
	}
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

var (
	encodingMap = map[string]string{ // all the other compress methods will ignore
		"gzip":     "gzip",
		"deflate":  "deflate",
		"*":        "gzip", // * means any compress will accept,we prefer gzip
		"identity": "",     // identity means none-compress
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
	return encodingMap[lastQ.name]
}
