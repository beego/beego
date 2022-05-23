// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"fmt"
	"path"
	"strconv"
)

var formatterMap = make(map[string]LogFormatter, 4)

type LogFormatter interface {
	Format(lm *LogMsg) string
}

// PatternLogFormatter provides a quick format method
// for example:
// tes := &PatternLogFormatter{Pattern: "%F:%n|%w %t>> %m", WhenFormat: "2006-01-02"}
// RegisterFormatter("tes", tes)
// SetGlobalFormatter("tes")
type PatternLogFormatter struct {
	Pattern    string
	WhenFormat string
}

func (p *PatternLogFormatter) getWhenFormatter() string {
	s := p.WhenFormat
	if s == "" {
		s = "2006/01/02 15:04:05.123" // default style
	}
	return s
}

func (p *PatternLogFormatter) Format(lm *LogMsg) string {
	return p.ToString(lm)
}

// RegisterFormatter register an formatter. Usually you should use this to extend your custom formatter
// for example:
// RegisterFormatter("my-fmt", &MyFormatter{})
// logs.SetFormatter(Console, `{"formatter": "my-fmt"}`)
func RegisterFormatter(name string, fmtr LogFormatter) {
	formatterMap[name] = fmtr
}

func GetFormatter(name string) (LogFormatter, bool) {
	res, ok := formatterMap[name]
	return res, ok
}

// ToString 'w' when, 'm' msg,'f' filename，'F' full path，'n' line number
// 'l' level number, 't' prefix of level type, 'T' full name of level type
func (p *PatternLogFormatter) ToString(lm *LogMsg) string {
	s := []rune(p.Pattern)
	msg := fmt.Sprintf(lm.Msg, lm.Args...)
	m := map[rune]string{
		'w': lm.When.Format(p.getWhenFormatter()),
		'm': msg,
		'n': strconv.Itoa(lm.LineNumber),
		'l': strconv.Itoa(lm.Level),
		't': levelPrefix[lm.Level],
		'T': levelNames[lm.Level],
		'F': lm.FilePath,
	}
	_, m['f'] = path.Split(lm.FilePath)
	res := ""
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '%' {
			if k, ok := m[s[i+1]]; ok {
				res += k
				i++
				continue
			}
		}
		res += string(s[i])
	}
	return res
}
