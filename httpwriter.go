// Copyright 2016 beego Author. All Rights Reserved.
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
	"bufio"
	"errors"
	"net"
	"net/http"
)

//response is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type responseWriter struct {
	rw          http.ResponseWriter
	started     bool
	status      int
	wroteHeader bool
}

func (r *responseWriter) reset(rw http.ResponseWriter) {
	r.rw = rw
	r.status = 0
	r.started = false
	r.wroteHeader = false
}



// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true.
// started means the response has sent out.
func (r *responseWriter) Write(p []byte) (int, error) {
	r.started = true
	if !r.wroteHeader {
		if r.status != 0 {
			r.rw.WriteHeader(r.status)
		} else {
			r.rw.WriteHeader(http.StatusOK)
		}
		r.wroteHeader = true
	}
	return r.rw.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (r *responseWriter) WriteHeader(code int) {
	if r.status > 0 {
		//prevent multiple response.WriteHeader calls
		return
	}
	r.status = code
	r.started = true
}

func (r *responseWriter) Header() http.Header {
	return r.rw.Header()
}

// Hijack hijacker for http
func (r *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.rw.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("webserver doesn't support hijacking")
	}
	return hj.Hijack()
}

// Flush http.Flusher
func (r *responseWriter) Flush() {
	if f, ok := r.rw.(http.Flusher); ok {
		f.Flush()
	}
}

// CloseNotify http.CloseNotifier
func (r *responseWriter) CloseNotify() <-chan bool {
	if cn, ok := r.rw.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}
