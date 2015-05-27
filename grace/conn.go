package grace

import "net"

type graceConn struct {
	net.Conn
	server *graceServer
}

func (c graceConn) Close() error {
	c.server.wg.Done()
	return c.Conn.Close()
}
