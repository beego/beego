// Zero-downtime restarts in Go.
package beego

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"strconv"
	"sync"
	//"syscall"
)

const (
	// An environment variable when restarting application http listener.
	FDKey = "BEEGO_HOT_FD"
)

// Export an error equivalent to net.errClosing for use with Accept during
// a graceful exit.
var ErrClosing = errors.New("use of closed network connection")
var ErrInitStart = errors.New("init from")

// Allows for us to notice when the connection is closed.
type conn struct {
	net.Conn
	wg      *sync.WaitGroup
	isclose bool
	lock    sync.Mutex
}

// Close current processing connection.
func (c conn) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.Conn.Close()
	if !c.isclose && err == nil {
		c.wg.Done()
		c.isclose = true
	}
	return err
}

type stoppableListener struct {
	net.Listener
	count   int64
	stopped bool
	wg      sync.WaitGroup
}

var theStoppable *stoppableListener

func newStoppable(l net.Listener) (sl *stoppableListener) {
	sl = &stoppableListener{Listener: l}

	// this goroutine monitors the channel. Can't do this in
	// Accept (below) because once it enters sl.Listener.Accept()
	// it blocks. We unblock it by closing the fd it is trying to
	// accept(2) on.
	go func() {
		WaitSignal(l)
		sl.stopped = true
		sl.Listener.Close()
	}()
	return
}

// Set stopped Listener to accept requests again.
// it returns the accepted and closable connection or error.
func (sl *stoppableListener) Accept() (c net.Conn, err error) {
	c, err = sl.Listener.Accept()
	if err != nil {
		return
	}
	sl.wg.Add(1)
	// Wrap the returned connection, so that we can observe when
	// it is closed.
	c = conn{Conn: c, wg: &sl.wg}

	return
}

// Listener waits signal to kill or interrupt then restart.
func WaitSignal(l net.Listener) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	for {
		sig := <-ch
		log.Println(sig.String())
		switch sig {

		case os.Kill:
			return nil
		case os.Interrupt:
			err := Restart(l)
			if nil != err {
				return err
			}
			return nil
		}
	}
	return nil // It'll never get here.
}

// Kill current running os process.
func CloseSelf() error {
	ppid := os.Getpid()
	if ppid == 1 { // init provided sockets, for example systemd
		return nil
	}
	p, err := os.FindProcess(ppid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// Re-exec this image without dropping the listener passed to this function.
func Restart(l net.Listener) error {
	argv0, err := exec.LookPath(os.Args[0])
	if nil != err {
		return err
	}
	wd, err := os.Getwd()
	if nil != err {
		return err
	}
	v := reflect.ValueOf(l).Elem().FieldByName("fd").Elem()
	fd := uintptr(v.FieldByName("sysfd").Int())
	allFiles := append([]*os.File{os.Stdin, os.Stdout, os.Stderr},
		os.NewFile(fd, string(v.FieldByName("sysfile").String())))

	p, err := os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   append(os.Environ(), fmt.Sprintf("%s=%d", FDKey, fd)),
		Files: allFiles,
	})
	if nil != err {
		return err
	}
	log.Printf("spawned child %d\n", p.Pid)
	return nil
}

// Get current net.Listen in running process.
func GetInitListener(tcpaddr *net.TCPAddr) (l net.Listener, err error) {
	countStr := os.Getenv(FDKey)
	if countStr == "" {
		return net.ListenTCP("tcp", tcpaddr)
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(count), "listen socket")
	l, err = net.FileListener(f)
	if err != nil {
		return nil, err
	}
	return l, nil
}
