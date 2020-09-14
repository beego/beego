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

var formatterMap = make(map[string]LogFormatter, 4)

type LogFormatter interface {
	Format(lm *LogMsg) string
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
