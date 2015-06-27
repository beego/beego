package grace

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type graceServer struct {
	*http.Server
	GraceListener    net.Listener
	SignalHooks      map[int]map[os.Signal][]func()
	tlsInnerListener *graceListener
	wg               sync.WaitGroup
	sigChan          chan os.Signal
	isChild          bool
	state            uint8
	Network          string
}

// Serve accepts incoming connections on the Listener l,
// creating a new service goroutine for each.
// The service goroutines read requests and then call srv.Handler to reply to them.
func (srv *graceServer) Serve() (err error) {
	srv.state = STATE_RUNNING
	err = srv.Server.Serve(srv.GraceListener)
	log.Println(syscall.Getpid(), "Waiting for connections to finish...")
	srv.wg.Wait()
	srv.state = STATE_TERMINATE
	return
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
// to handle requests on incoming connections. If srv.Addr is blank, ":http" is
// used.
func (srv *graceServer) ListenAndServe() (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	go srv.handleSignals()

	l, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	srv.GraceListener = newGraceListener(l, srv)

	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Kill()
		if err != nil {
			return err
		}
	}

	log.Println(os.Getpid(), srv.Addr)
	return srv.Serve()
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for the server must
// be provided. If the certificate is signed by a certificate authority, the
// certFile should be the concatenation of the server's certificate followed by the
// CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
func (srv *graceServer) ListenAndServeTLS(certFile, keyFile string) (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	go srv.handleSignals()

	l, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	srv.tlsInnerListener = newGraceListener(l, srv)
	srv.GraceListener = tls.NewListener(srv.tlsInnerListener, config)

	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Kill()
		if err != nil {
			return err
		}
	}
	log.Println(os.Getpid(), srv.Addr)
	return srv.Serve()
}

// getListener either opens a new socket to listen on, or takes the acceptor socket
// it got passed when restarted.
func (srv *graceServer) getListener(laddr string) (l net.Listener, err error) {
	if srv.isChild {
		var ptrOffset uint = 0
		if len(socketPtrOffsetMap) > 0 {
			ptrOffset = socketPtrOffsetMap[laddr]
			log.Println("laddr", laddr, "ptr offset", socketPtrOffsetMap[laddr])
		}

		f := os.NewFile(uintptr(3+ptrOffset), "")
		l, err = net.FileListener(f)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return
		}
	} else {
		l, err = net.Listen(srv.Network, laddr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return
		}
	}
	return
}

// handleSignals listens for os Signals and calls any hooked in function that the
// user had registered with the signal.
func (srv *graceServer) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	pid := syscall.Getpid()
	for {
		sig = <-srv.sigChan
		srv.signalHooks(PRE_SIGNAL, sig)
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "Received SIGHUP. forking.")
			err := srv.fork()
			if err != nil {
				log.Println("Fork err:", err)
			}
		case syscall.SIGINT:
			log.Println(pid, "Received SIGINT.")
			srv.shutdown()
		case syscall.SIGTERM:
			log.Println(pid, "Received SIGTERM.")
			srv.shutdown()
		default:
			log.Printf("Received %v: nothing i care about...\n", sig)
		}
		srv.signalHooks(POST_SIGNAL, sig)
	}
}

func (srv *graceServer) signalHooks(ppFlag int, sig os.Signal) {
	if _, notSet := srv.SignalHooks[ppFlag][sig]; !notSet {
		return
	}
	for _, f := range srv.SignalHooks[ppFlag][sig] {
		f()
	}
	return
}

// shutdown closes the listener so that no new connections are accepted. it also
// starts a goroutine that will serverTimeout (stop all running requests) the server
// after DefaultTimeout.
func (srv *graceServer) shutdown() {
	if srv.state != STATE_RUNNING {
		return
	}

	srv.state = STATE_SHUTTING_DOWN
	if DefaultTimeout >= 0 {
		go srv.serverTimeout(DefaultTimeout)
	}
	err := srv.GraceListener.Close()
	if err != nil {
		log.Println(syscall.Getpid(), "Listener.Close() error:", err)
	} else {
		log.Println(syscall.Getpid(), srv.GraceListener.Addr(), "Listener closed.")
	}
}

// serverTimeout forces the server to shutdown in a given timeout - whether it
// finished outstanding requests or not. if Read/WriteTimeout are not set or the
// max header size is very big a connection could hang
func (srv *graceServer) serverTimeout(d time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("WaitGroup at 0", r)
		}
	}()
	if srv.state != STATE_SHUTTING_DOWN {
		return
	}
	time.Sleep(d)
	log.Println("[STOP - Hammer Time] Forcefully shutting down parent")
	for {
		if srv.state == STATE_TERMINATE {
			break
		}
		srv.wg.Done()
	}
}

func (srv *graceServer) fork() (err error) {
	regLock.Lock()
	defer regLock.Unlock()
	if runningServersForked {
		return
	}
	runningServersForked = true

	var files = make([]*os.File, len(runningServers))
	var orderArgs = make([]string, len(runningServers))
	for _, srvPtr := range runningServers {
		switch srvPtr.GraceListener.(type) {
		case *graceListener:
			files[socketPtrOffsetMap[srvPtr.Server.Addr]] = srvPtr.GraceListener.(*graceListener).File()
		default:
			files[socketPtrOffsetMap[srvPtr.Server.Addr]] = srvPtr.tlsInnerListener.File()
		}
		orderArgs[socketPtrOffsetMap[srvPtr.Server.Addr]] = srvPtr.Server.Addr
	}

	log.Println(files)
	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-graceful" {
				break
			}
			args = append(args, arg)
		}
	}
	args = append(args, "-graceful")
	if len(runningServers) > 1 {
		args = append(args, fmt.Sprintf(`-socketorder=%s`, strings.Join(orderArgs, ",")))
		log.Println(args)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Restart: Failed to launch, error: %v", err)
	}

	return
}
