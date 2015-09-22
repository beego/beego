package grace

import "net"

type graceConn struct {
	net.Conn
	server *Server
}

func (c graceConn) Close() error {
	c.server.wg.Done()
	return c.Conn.Close()
}
