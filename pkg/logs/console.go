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
	"os"
	"strings"

	"github.com/shiena/ansicolor"
)

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
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
	OldLoggerAdapter
	lg       *logWriter
	fmtter LogFormatter
	Level    int  `json:"level"`
	Colorful bool `json:"color"` //this filed is useful only when system's terminal supports color
}

// NewConsole create ConsoleWriter returning as LoggerInterface.
func NewConsole() Logger {
	cw := &consoleWriter{
		lg:       newLogWriter(ansicolor.NewAnsiColorWriter(os.Stdout)),
		Level:    LevelDebug,
		Colorful: true,
		fmtter: &consoleDefaultFormatter {
			colorful: true,
		},
	}
	return cw
}

// Init init console logger.
// jsonConfig like '{"level":LevelTrace}'.
func (c *consoleWriter) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	return json.Unmarshal([]byte(jsonConfig), c)
}

type consoleDefaultFormatter struct {
	colorful bool
}

func (cdf *consoleDefaultFormatter) Format(lm *LogMsg) string {
	msg := lm.msg
	if cdf.colorful {
		msg = strings.Replace(lm.msg, levelPrefix[lm.level], colors[lm.level](levelPrefix[lm.level]), 1)
	}

	h, _, _ := formatTimeHeader(lm.when)

	bytes := append(append(h, msg...), '\n')

	return string(bytes)
}

// WriteMsg write message in console.
func (c *consoleWriter) WriteLogMsg(lm *LogMsg) error {

	// here is an example

	if lm.level > c.Level {
		return nil
	}

	msg := c.fmtter.Format(lm)
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
