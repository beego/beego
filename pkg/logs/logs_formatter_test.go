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
	"log"
	"os"
	"testing"

	"github.com/astaxie/beego/pkg/common"
)

type customLogFormatter struct {
}

func (clf *customLogFormatter) Format(lm *LogMsg) string {
	return "[CONSOLE]" + lm.FilePath + lm.When.String() + lm.Msg
}

func (clf *customLogFormatter) Format2(lm *LogMsg) string {
	return "[FILE]" + lm.FilePath + lm.When.String() + lm.Msg
}

func (clf *customLogFormatter) GlobalFormat(lm *LogMsg) string {
	return "[Global]" + lm.FilePath + lm.When.String() + lm.Msg
}

func TestCustomFormatter(t *testing.T) {
	fmtter := &customLogFormatter{}

	SetLoggerWithOpts("console", []string{""}, common.SimpleKV{Key: "formatter", Value: fmtter.Format})
	return
}

func TestGlobalFormatter(t *testing.T) {
	fmtter := &customLogFormatter{}

	err := SetLoggerWithOpts("console", []string{""}, common.SimpleKV{Key: "formatter", Value: fmtter.Format})
	if err != nil {
		log.Fatal(err)
	}

	err = SetLoggerWithOpts("file", []string{`{"filename":"tester.log"}`}, common.SimpleKV{Key: "formatter", Value: fmtter.Format2})
	if err != nil {
		log.Fatal(err)
	}

	os.Remove("tester.log")

	SetGlobalFormatter(fmtter.GlobalFormat)
}

func TestDuplicateAdapterFormatter(t *testing.T) {
	fmtter := &customLogFormatter{}

	SetLogger("console")

	err := SetLoggerWithOpts("console", []string{``}, common.SimpleKV{Key: "formatter", Value: fmtter.Format2})
	if err == nil {
		t.Fatal("duplicate log adapter console was set without failure")
	}

}
