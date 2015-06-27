package grace

import (
	"net"
	"os"
	"syscall"
	"time"
)

type graceListener struct {
	net.Listener
	stop    chan error
	stopped bool
	server  *graceServer
}

func newGraceListener(l net.Listener, srv *graceServer) (el *graceListener) {
	el = &graceListener{
		Listener: l,
		stop:     make(chan error),
		server:   srv,
	}
	go func() {
		_ = <-el.stop
		el.stopped = true
		el.stop <- el.Listener.Close()
	}()
	return
}

func (gl *graceListener) Accept() (c net.Conn, err error) {
	tc, err := gl.Listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)

	c = graceConn{
		Conn:   tc,
		server: gl.server,
	}

	gl.server.wg.Add(1)
	return
}

func (el *graceListener) Close() error {
	if el.stopped {
		return syscall.EINVAL
	}
	el.stop <- nil
	return <-el.stop
}

func (el *graceListener) File() *os.File {
	// returns a dup(2) - FD_CLOEXEC flag *not* set
	tl := el.Listener.(*net.TCPListener)
	fl, _ := tl.File()
	return fl
}
