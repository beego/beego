package grace

import (
	"os"

	"github.com/beego/beego/server/web/grace"
)

// Server embedded http.Server
type Server grace.Server

// Serve accepts incoming connections on the Listener l,
// creating a new service goroutine for each.
// The service goroutines read requests and then call srv.Handler to reply to them.
func (srv *Server) Serve() (err error) {
	return (*grace.Server)(srv).Serve()
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
// to handle requests on incoming connections. If srv.Addr is blank, ":http" is
// used.
func (srv *Server) ListenAndServe() (err error) {
	return (*grace.Server)(srv).ListenAndServe()
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
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	return (*grace.Server)(srv).ListenAndServeTLS(certFile, keyFile)
}

// ListenAndServeMutualTLS listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming mutual TLS connections.
func (srv *Server) ListenAndServeMutualTLS(certFile, keyFile, trustFile string) error {
	return (*grace.Server)(srv).ListenAndServeMutualTLS(certFile, keyFile, trustFile)
}

// RegisterSignalHook registers a function to be run PreSignal or PostSignal for a given signal.
func (srv *Server) RegisterSignalHook(ppFlag int, sig os.Signal, f func()) error {
	return (*grace.Server)(srv).RegisterSignalHook(ppFlag, sig, f)
}
