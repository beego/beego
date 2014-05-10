// +build windows
// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package logs

import "syscall"

var colors = [6]int{
	3, // Trace      cyan or Aqua
	1, // Debug      blue
	2, // Info       green
	6, // Warn       yellow
	4, // Error      red
	5, // Critical   purple
}


// write message in console.
func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level < c.Level {
		return nil
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(colors[level]))
	c.lg.Println(msg)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
	return nil
}
