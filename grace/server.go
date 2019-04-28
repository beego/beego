package grace

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Server embedded http.Server
type Server struct {
	*http.Server
	ln           net.Listener
	SignalHooks  map[int]map[os.Signal][]func()
	sigChan      chan os.Signal
	isChild      bool
	state        uint8
	Network      string
	terminalChan chan error
}

// Serve accepts incoming connections on the Listener l,
// creating a new service goroutine for each.
// The service goroutines read requests and then call srv.Handler to reply to them.
func (srv *Server) Serve() (err error) {
	srv.state = StateRunning
	defer func() { srv.state = StateTerminate }()

	// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS
	// immediately return ErrServerClosed. Make sure the program doesn't exit
	// and waits instead for Shutdown to return.
	if err = srv.Server.Serve(srv.ln); err != nil && err != http.ErrServerClosed {
		log.Println(syscall.Getpid(), "Server.Serve() error:", err)
		return err
	}

	log.Println(syscall.Getpid(), srv.ln.Addr(), "Listener closed.")
	// wait for Shutdown to return
	return <-srv.terminalChan
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
// to handle requests on incoming connections. If srv.Addr is blank, ":http" is
// used.
func (srv *Server) ListenAndServe() (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	go srv.handleSignals()

	srv.ln, err = srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Signal(syscall.SIGTERM)
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
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	if srv.TLSConfig == nil {
		srv.TLSConfig = &tls.Config{}
	}
	if srv.TLSConfig.NextProtos == nil {
		srv.TLSConfig.NextProtos = []string{"http/1.1"}
	}

	srv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	srv.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	go srv.handleSignals()

	ln, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}
	srv.ln = tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)

	if srv.isChild {
		process, err := os.FindProcess(os.Getppid())
		if err != nil {
			log.Println(err)
			return err
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}

	log.Println(os.Getpid(), srv.Addr)
	return srv.Serve()
}

// ListenAndServeMutualTLS listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming mutual TLS connections.
func (srv *Server) ListenAndServeMutualTLS(certFile, keyFile, trustFile string) (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	if srv.TLSConfig == nil {
		srv.TLSConfig = &tls.Config{}
	}
	if srv.TLSConfig.NextProtos == nil {
		srv.TLSConfig.NextProtos = []string{"http/1.1"}
	}

	srv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	srv.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	srv.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	pool := x509.NewCertPool()
	data, err := ioutil.ReadFile(trustFile)
	if err != nil {
		log.Println(err)
		return err
	}
	pool.AppendCertsFromPEM(data)
	srv.TLSConfig.ClientCAs = pool
	log.Println("Mutual HTTPS")
	go srv.handleSignals()

	ln, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}
	srv.ln = tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, srv.TLSConfig)

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
func (srv *Server) getListener(laddr string) (l net.Listener, err error) {
	if srv.isChild {
		var ptrOffset uint
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

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// handleSignals listens for os Signals and calls any hooked in function that the
// user had registered with the signal.
func (srv *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.sigChan,
		hookableSignals...,
	)

	pid := syscall.Getpid()
	for {
		sig = <-srv.sigChan
		srv.signalHooks(PreSignal, sig)
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
		srv.signalHooks(PostSignal, sig)
	}
}

func (srv *Server) signalHooks(ppFlag int, sig os.Signal) {
	if _, notSet := srv.SignalHooks[ppFlag][sig]; !notSet {
		return
	}
	for _, f := range srv.SignalHooks[ppFlag][sig] {
		f()
	}
}

// shutdown closes the listener so that no new connections are accepted. it also
// starts a goroutine that will serverTimeout (stop all running requests) the server
// after DefaultTimeout.
func (srv *Server) shutdown() {
	if srv.state != StateRunning {
		return
	}

	srv.state = StateShuttingDown
	log.Println(syscall.Getpid(), "Waiting for connections to finish...")
	ctx := context.Background()
	if DefaultTimeout >= 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
	}
	srv.terminalChan <- srv.Server.Shutdown(ctx)
}

func (srv *Server) fork() (err error) {
	regLock.Lock()
	defer regLock.Unlock()
	if runningServersForked {
		return
	}
	runningServersForked = true

	var files = make([]*os.File, len(runningServers))
	var orderArgs = make([]string, len(runningServers))
	for _, srvPtr := range runningServers {
		f, _ := srvPtr.ln.(*net.TCPListener).File()
		files[socketPtrOffsetMap[srvPtr.Server.Addr]] = f
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

// RegisterSignalHook registers a function to be run PreSignal or PostSignal for a given signal.
func (srv *Server) RegisterSignalHook(ppFlag int, sig os.Signal, f func()) (err error) {
	if ppFlag != PreSignal && ppFlag != PostSignal {
		err = fmt.Errorf("Invalid ppFlag argument. Must be either grace.PreSignal or grace.PostSignal")
		return
	}
	for _, s := range hookableSignals {
		if s == sig {
			srv.SignalHooks[ppFlag][sig] = append(srv.SignalHooks[ppFlag][sig], f)
			return
		}
	}
	err = fmt.Errorf("Signal '%v' is not supported", sig)
	return
}
