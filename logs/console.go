package logs

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleWriter struct {
	lg    *log.Logger
	level int
}

func NewConsole() LoggerInterface {
	cw := new(ConsoleWriter)
	cw.lg = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cw.level = LevelTrace
	return cw
}

func (c *ConsoleWriter) Init(jsonconfig string) error {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonconfig), &m)
	if err != nil {
		return err
	}
	if lv, ok := m["level"]; ok {
		c.level = int(lv.(float64))
	}
	return nil
}

func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level < c.level {
		return nil
	}
	c.lg.Println(msg)
	return nil
}

func (c *ConsoleWriter) Destroy() {

}

func (c *ConsoleWriter) Flush() {

}

func init() {
	Register("console", NewConsole)
}
