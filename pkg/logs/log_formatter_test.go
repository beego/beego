// Copyright 2020 beego 
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
	"github.com/astaxie/beego/pkg/common"
)

type customLogFormatter struct {

}

func (clf *customLogFormatter) Format(lm *LogMsg) string {
	return lm.filePath + lm.when.String() + lm.msg
}

func Example()  {
	fmtter := &customLogFormatter{}
	SetGlobalFormatter(fmtter)
	// or
	SetLoggerWithOpts(AdapterAliLS, []string{""}, common.SimpleKV{ Key:"formatter", Value:fmtter})

	// to resole
}
