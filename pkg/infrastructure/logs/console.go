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

	"github.com/astaxie/beego/pkg/infrastructure/utils"

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
	lg              *logWriter
	customFormatter func(*LogMsg) string
	Level           int  `json:"level"`
	Colorful        bool `json:"color"` //this filed is useful only when system's terminal supports color
}

func (c *consoleWriter) Format(lm *LogMsg) string {
	msg := lm.Msg

	h, _, _ := formatTimeHeader(lm.When)
	bytes := append(append(h, msg...), '\n')

	return string(bytes)

}

// NewConsole creates ConsoleWriter returning as LoggerInterface.
func NewConsole() Logger {
	cw := &consoleWriter{
		lg:       newLogWriter(ansicolor.NewAnsiColorWriter(os.Stdout)),
		Level:    LevelDebug,
		Colorful: true,
	}
	return cw
}

// Init initianlizes the console logger.
// jsonConfig must be in the format '{"level":LevelTrace}'
func (c *consoleWriter) Init(jsonConfig string, opts ...utils.KV) error {

	for _, elem := range opts {
		if elem.GetKey() == "formatter" {
			formatter, err := GetFormatter(elem)
			if err != nil {
				return err
			}
			c.customFormatter = formatter
		}
	}

	if len(jsonConfig) == 0 {
		return nil
	}

	return json.Unmarshal([]byte(jsonConfig), c)
}

// WriteMsg writes message in console.
func (c *consoleWriter) WriteMsg(lm *LogMsg) error {
	if lm.Level > c.Level {
		return nil
	}

	msg := ""

	if c.Colorful {
		lm.Msg = strings.Replace(lm.Msg, levelPrefix[lm.Level], colors[lm.Level](levelPrefix[lm.Level]), 1)
	}

	if c.customFormatter != nil {
		msg = c.customFormatter(lm)
	} else {
		msg = c.Format(lm)

	}

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
