// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package logs

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"syscall"
)

type Brush func(string) string

func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;36"), // Trace      cyan
	NewBrush("1;34"), // Debug      blue
	NewBrush("1;32"), // Info       green
	NewBrush("1;33"), // Warn       yellow
	NewBrush("1;31"), // Error      red
	NewBrush("1;35"), // Critical   purple
}
var windowsColors = []int{
	3, // Trace      cyan or Aqua
	1, // Debug      blue
	2, // Info       green
	6, // Warn       yellow
	4, // Error      red
	5, // Critical   purple
}

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

// write message in console.
func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level < c.Level {
		return nil
	}
	if goos := runtime.GOOS; goos == "windows" {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		proc := kernel32.NewProc("SetConsoleTextAttribute")
		handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(windowsColors[level]))
		c.lg.Println(msg)
		handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
		CloseHandle := kernel32.NewProc("CloseHandle")
		CloseHandle.Call(handle)
	} else {
		c.lg.Println(colors[level](msg))
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
