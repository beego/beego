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
	"io"
	"net/http"
	"os"

	"github.com/astaxie/beego/pkg/server/web/context"
)

// InitGzip init the gzipcompress
func InitGzip(minLength, compressLevel int, methods []string) {
	context.InitGzip(minLength, compressLevel, methods)
}

// WriteFile reads from file and writes to writer by the specific encoding(gzip/deflate)
func WriteFile(encoding string, writer io.Writer, file *os.File) (bool, string, error) {
	return context.WriteFile(encoding, writer, file)
}

// WriteBody reads  writes content to writer by the specific encoding(gzip/deflate)
func WriteBody(encoding string, writer io.Writer, content []byte) (bool, string, error) {
	return context.WriteBody(encoding, writer, content)
}

// ParseEncoding will extract the right encoding for response
// the Accept-Encoding's sec is here:
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.3
func ParseEncoding(r *http.Request) string {
	return context.ParseEncoding(r)
}
