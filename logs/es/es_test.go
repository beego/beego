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

package es

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

// Try each log level in decreasing order of priority.
func testESCalls(bl *logs.BeeLogger) {
	bl.Emergency("emergency")
	bl.Alert("alert")
	bl.Critical("critical")
	bl.Error("error")
	bl.Warning("warning")
	bl.Notice("notice")
	bl.Informational("informational")
	bl.Debug("debug")
}

func TestLogToES(t *testing.T) {
	log := logs.NewLogger(100)
	err := log.SetLogger("es", `{"dsn":"http://localhost:9200/","index_prefix":"test-name-","index":"test-index-name","level":7}`)
	if err != nil {
		t.Errorf("set logger adapter error")
	}
	testESCalls(log)
}
