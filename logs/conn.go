package logs

import (
	"encoding/json"
	"io"
	"log"
	"net"
)

type ConnWriter struct {
	lg             *log.Logger
	innerWriter    io.WriteCloser
	reconnectOnMsg bool
	reconnect      bool
	net            string
	addr           string
	level          int
}

func NewConn() LoggerInterface {
	conn := new(ConnWriter)
	conn.level = LevelTrace
	return conn
}

func (c *ConnWriter) Init(jsonconfig string) error {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonconfig), &m)
	if err != nil {
		return err
	}
	if rom, ok := m["reconnectOnMsg"]; ok {
		c.reconnectOnMsg = rom.(bool)
	}
	if rc, ok := m["reconnect"]; ok {
		c.reconnect = rc.(bool)
	}
	if nt, ok := m["net"]; ok {
		c.net = nt.(string)
	}
	if addr, ok := m["addr"]; ok {
		c.addr = addr.(string)
	}
	if lv, ok := m["level"]; ok {
		c.level = int(lv.(float64))
	}
	return nil
}

func (c *ConnWriter) WriteMsg(msg string, level int) error {
	if level < c.level {
		return nil
	}
	if c.neddedConnectOnMsg() {
		err := c.connect()
		if err != nil {
			return err
		}
	}

	if c.reconnectOnMsg {
		defer c.innerWriter.Close()
	}
	c.lg.Println(msg)
	return nil
}

func (c *ConnWriter) Flush() {

}

func (c *ConnWriter) Destroy() {
	if c.innerWriter == nil {
		return
	}
	c.innerWriter.Close()
}

func (c *ConnWriter) connect() error {
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}

	conn, err := net.Dial(c.net, c.addr)
	if err != nil {
		return err
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		tcpConn.SetKeepAlive(true)
	}

	c.innerWriter = conn
	c.lg = log.New(conn, "", log.Ldate|log.Ltime)
	return nil
}

func (c *ConnWriter) neddedConnectOnMsg() bool {
	if c.reconnect {
		c.reconnect = false
		return true
	}

	if c.innerWriter == nil {
		return true
	}

	return c.reconnectOnMsg
}

func init() {
	Register("conn", NewConn)
}
