// +build !windows
// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/astaxie/beego for the canonical source repository
// @license     http://github.com/astaxie/beego/blob/master/LICENSE
// @authors     astaxie

package logs

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

// write message in console.
func (c *ConsoleWriter) WriteMsg(msg string, level int) error {
	if level < c.Level {
		return nil
	}
	c.lg.Println(colors[level](msg))
	return nil
}
