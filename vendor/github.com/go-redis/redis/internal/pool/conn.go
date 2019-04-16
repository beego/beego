package pool

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/internal/proto"
)

var noDeadline = time.Time{}

type Conn struct {
	netConn net.Conn

	Rd *proto.Reader
	wb *proto.WriteBuffer

	concurrentReadWrite bool

	Inited bool
	usedAt atomic.Value
}

func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		netConn: netConn,
	}
	buf := proto.NewElasticBufReader(netConn)
	cn.Rd = proto.NewReader(buf)
	cn.wb = proto.NewWriteBuffer()
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *Conn) UsedAt() time.Time {
	return cn.usedAt.Load().(time.Time)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	cn.usedAt.Store(tm)
}

func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
	cn.Rd.Reset(netConn)
}

func (cn *Conn) IsStale(timeout time.Duration) bool {
	return timeout > 0 && time.Since(cn.UsedAt()) > timeout
}

func (cn *Conn) SetReadTimeout(timeout time.Duration) {
	now := time.Now()
	cn.SetUsedAt(now)
	if timeout > 0 {
		cn.netConn.SetReadDeadline(now.Add(timeout))
	} else {
		cn.netConn.SetReadDeadline(noDeadline)
	}
}

func (cn *Conn) SetWriteTimeout(timeout time.Duration) {
	now := time.Now()
	cn.SetUsedAt(now)
	if timeout > 0 {
		cn.netConn.SetWriteDeadline(now.Add(timeout))
	} else {
		cn.netConn.SetWriteDeadline(noDeadline)
	}
}

func (cn *Conn) Write(b []byte) (int, error) {
	return cn.netConn.Write(b)
}

func (cn *Conn) RemoteAddr() net.Addr {
	return cn.netConn.RemoteAddr()
}

func (cn *Conn) EnableConcurrentReadWrite() {
	cn.concurrentReadWrite = true
	cn.wb.ResetBuffer(make([]byte, 4096))
}

func (cn *Conn) PrepareWriteBuffer() *proto.WriteBuffer {
	if !cn.concurrentReadWrite {
		cn.wb.ResetBuffer(cn.Rd.Buffer())
	}
	cn.wb.Reset()
	return cn.wb
}

func (cn *Conn) FlushWriteBuffer(wb *proto.WriteBuffer) error {
	_, err := cn.netConn.Write(wb.Bytes())
	if !cn.concurrentReadWrite {
		cn.Rd.ResetBuffer(wb.Buffer())
	}
	return err
}

func (cn *Conn) Close() error {
	return cn.netConn.Close()
}
