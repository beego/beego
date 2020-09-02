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

// Package grace use to hot reload
// Description: http://grisha.org/blog/2014/06/03/graceful-restart-in-golang/
//
// Usage:
//
// import(
//   "log"
//	 "net/http"
//	 "os"
//
//   "github.com/astaxie/beego/grace"
// )
//
//  func handler(w http.ResponseWriter, r *http.Request) {
//	  w.Write([]byte("WORLD!"))
//  }
//
//  func main() {
//      mux := http.NewServeMux()
//      mux.HandleFunc("/hello", handler)
//
//	    err := grace.ListenAndServe("localhost:8080", mux)
//      if err != nil {
//		   log.Println(err)
//	    }
//      log.Println("Server on 8080 stopped")
//	     os.Exit(0)
//    }
package grace

import (
	"net/http"
	"time"

	"github.com/astaxie/beego/pkg/server/web/grace"
)

const (
	// PreSignal is the position to add filter before signal
	PreSignal = iota
	// PostSignal is the position to add filter after signal
	PostSignal
	// StateInit represent the application inited
	StateInit
	// StateRunning represent the application is running
	StateRunning
	// StateShuttingDown represent the application is shutting down
	StateShuttingDown
	// StateTerminate represent the application is killed
	StateTerminate
)

var (


	// DefaultReadTimeOut is the HTTP read timeout
	DefaultReadTimeOut time.Duration
	// DefaultWriteTimeOut is the HTTP Write timeout
	DefaultWriteTimeOut time.Duration
	// DefaultMaxHeaderBytes is the Max HTTP Header size, default is 0, no limit
	DefaultMaxHeaderBytes int
	// DefaultTimeout is the shutdown server's timeout. default is 60s
	DefaultTimeout = grace.DefaultTimeout

)

// NewServer returns a new graceServer.
func NewServer(addr string, handler http.Handler) (srv *Server) {
	return (*Server)(grace.NewServer(addr, handler))
}

// ListenAndServe refer http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

// ListenAndServeTLS refer http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}
