// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/shiena/ansicolor"
)

// brush is a color join function
type brush func(string) string

// newBrush returns a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;37"), // Emergency          white
	newBrush("1;36"), // Alert              cyan
	newBrush("1;35"), // Critical           magenta
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;32"), // Notice             green
	newBrush("1;34"), // Informational      blue
	newBrush("1;44"), // Debug              Background blue
}

// consoleWriter implements LoggerInterface and writes messages to terminal.
type consoleWriter struct {
	lg        *logWriter
	formatter LogFormatter
	Formatter string `json:"formatter"`
	Level     int    `json:"level"`
	Colorful  bool   `json:"color"` // this filed is useful only when system's terminal supports color
}

func (c *consoleWriter) Format(lm *LogMsg) string {
	msg := lm.OldStyleFormat()
	if c.Colorful {
		msg = strings.Replace(msg, levelPrefix[lm.Level], colors[lm.Level](levelPrefix[lm.Level]), 1)
	}
	h, _, _ := formatTimeHeader(lm.When)
	return string(append(h, msg...))
}

func (c *consoleWriter) SetFormatter(f LogFormatter) {
	c.formatter = f
}

// NewConsole creates ConsoleWriter returning as LoggerInterface.
func NewConsole() Logger {
	return newConsole()
}

func newConsole() *consoleWriter {
	cw := &consoleWriter{
		lg:       newLogWriter(ansicolor.NewAnsiColorWriter(os.Stdout)),
		Level:    LevelDebug,
		Colorful: true,
	}
	cw.formatter = cw
	return cw
}

// Init initianlizes the console logger.
// jsonConfig must be in the format '{"level":LevelTrace}'
func (c *consoleWriter) Init(config string) error {
	if len(config) == 0 {
		return nil
	}

	res := json.Unmarshal([]byte(config), c)
	if res == nil && len(c.Formatter) > 0 {
		fmtr, ok := GetFormatter(c.Formatter)
		if !ok {
			return fmt.Errorf("the formatter with name: %s not found", c.Formatter)
		}
		c.formatter = fmtr
	}
	return res
}

// WriteMsg writes message in console.
func (c *consoleWriter) WriteMsg(lm *LogMsg) error {
	if lm.Level > c.Level {
		return nil
	}
	msg := c.formatter.Format(lm)
	c.lg.writeln(msg)
	return nil
}

// Destroy implementing method. empty.
func (c *consoleWriter) Destroy() {
}

// Flush implementing method. empty.
func (c *consoleWriter) Flush() {
}

func init() {
	Register(AdapterConsole, NewConsole)
}
