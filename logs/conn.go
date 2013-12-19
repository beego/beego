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
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          int    `json:"level"`
}

func NewConn() LoggerInterface {
	conn := new(ConnWriter)
	conn.Level = LevelTrace
	return conn
}

func (c *ConnWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConnWriter) WriteMsg(msg string, level int) error {
	if level < c.Level {
		return nil
	}
	if c.neddedConnectOnMsg() {
		err := c.connect()
		if err != nil {
			return err
		}
	}

	if c.ReconnectOnMsg {
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

	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
	}

	c.innerWriter = conn
	c.lg = log.New(conn, "", log.Ldate|log.Ltime)
	return nil
}

func (c *ConnWriter) neddedConnectOnMsg() bool {
	if c.Reconnect {
		c.Reconnect = false
		return true
	}

	if c.innerWriter == nil {
		return true
	}

	return c.ReconnectOnMsg
}

func init() {
	Register("conn", NewConn)
}
