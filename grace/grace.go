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
//	    err := grace.ListenAndServe("localhost:8080", mux1)
//      if err != nil {
//		   log.Println(err)
//	    }
//      log.Println("Server on 8080 stopped")
//	     os.Exit(0)
//    }
package grace

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	PRE_SIGNAL = iota
	POST_SIGNAL

	STATE_INIT
	STATE_RUNNING
	STATE_SHUTTING_DOWN
	STATE_TERMINATE
)

var (
	regLock              *sync.Mutex
	runningServers       map[string]*graceServer
	runningServersOrder  []string
	socketPtrOffsetMap   map[string]uint
	runningServersForked bool

	DefaultReadTimeOut    time.Duration
	DefaultWriteTimeOut   time.Duration
	DefaultMaxHeaderBytes int
	DefaultTimeout        time.Duration

	isChild     bool
	socketOrder string
)

func init() {
	regLock = &sync.Mutex{}
	flag.BoolVar(&isChild, "graceful", false, "listen on open fd (after forking)")
	flag.StringVar(&socketOrder, "socketorder", "", "previous initialization order - used when more than one listener was started")
	runningServers = make(map[string]*graceServer)
	runningServersOrder = []string{}
	socketPtrOffsetMap = make(map[string]uint)

	DefaultMaxHeaderBytes = 0

	DefaultTimeout = 60 * time.Second
}

// NewServer returns a new graceServer.
func NewServer(addr string, handler http.Handler) (srv *graceServer) {
	regLock.Lock()
	defer regLock.Unlock()
	if !flag.Parsed() {
		flag.Parse()
	}
	if len(socketOrder) > 0 {
		for i, addr := range strings.Split(socketOrder, ",") {
			socketPtrOffsetMap[addr] = uint(i)
		}
	} else {
		socketPtrOffsetMap[addr] = uint(len(runningServersOrder))
	}

	srv = &graceServer{
		wg:      sync.WaitGroup{},
		sigChan: make(chan os.Signal),
		isChild: isChild,
		SignalHooks: map[int]map[os.Signal][]func(){
			PRE_SIGNAL: map[os.Signal][]func(){
				syscall.SIGHUP:  []func(){},
				syscall.SIGINT:  []func(){},
				syscall.SIGTERM: []func(){},
			},
			POST_SIGNAL: map[os.Signal][]func(){
				syscall.SIGHUP:  []func(){},
				syscall.SIGINT:  []func(){},
				syscall.SIGTERM: []func(){},
			},
		},
		state:   STATE_INIT,
		Network: "tcp",
	}
	srv.Server = &http.Server{}
	srv.Server.Addr = addr
	srv.Server.ReadTimeout = DefaultReadTimeOut
	srv.Server.WriteTimeout = DefaultWriteTimeOut
	srv.Server.MaxHeaderBytes = DefaultMaxHeaderBytes
	srv.Server.Handler = handler

	runningServersOrder = append(runningServersOrder, addr)
	runningServers[addr] = srv

	return
}

// refer http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

// refer http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}
