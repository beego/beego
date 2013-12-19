package logs

import (
	"encoding/json"
	"log"
	"os"
)

type ConsoleWriter struct {
	lg    *log.Logger
	Level int `json:"level"`
}

func NewConsole() LoggerInterface {
	cw := new(ConsoleWriter)
	cw.lg = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cw.Level = LevelTrace
	return cw
}

func (c *ConsoleWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level < c.Level {
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
