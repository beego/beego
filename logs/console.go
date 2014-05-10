// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package logs

import (
	"encoding/json"
	"log"
	"os")

// ConsoleWriter implements LoggerInterface and writes messages to terminal.
type ConsoleWriter struct {
	lg    *log.Logger
	Level int `json:"level"`
}

// create ConsoleWriter returning as LoggerInterface.
func NewConsole() LoggerInterface {
	cw := new(ConsoleWriter)
	cw.lg = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	cw.Level = LevelTrace
	return cw
}

// init console logger.
// jsonconfig like '{"level":LevelTrace}'.
func (c *ConsoleWriter) Init(jsonconfig string) error {
	if len(jsonconfig) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(jsonconfig), c)
	if err != nil {
		return err
	}
	return nil
}



// implementing method. empty.
func (c *ConsoleWriter) Destroy() {

}

// implementing method. empty.
func (c *ConsoleWriter) Flush() {

}

func init() {
	Register("console", NewConsole)
}
