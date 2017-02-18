package grace

import (
	"errors"
	"log"
	"net"
)

type graceConn struct {
	net.Conn
	server *Server
}

func (c graceConn) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	err = c.Conn.Close()
	if err == nil {
		c.server.wg.Done()
	} else {
		log.Panicln("close error:", err)
	}
	return err
}
