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

package beego

import (
	"strings"

	"github.com/astaxie/beego/logs"
)

// Log levels to control the logging output.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	BeeLogger.SetLevel(l)
}

func SetLogFuncCall(b bool) {
	BeeLogger.EnableFuncCallDepth(b)
	BeeLogger.SetLogFuncCallDepth(3)
}

// logger references the used application logger.
var BeeLogger *logs.BeeLogger

// SetLogger sets a new logger.
func SetLogger(adaptername string, config string) error {
	err := BeeLogger.SetLogger(adaptername, config)
	if err != nil {
		return err
	}
	return nil
}

func Emergency(v ...interface{}) {
	BeeLogger.Emergency(generateFmtStr(len(v)), v...)
}

func Alert(v ...interface{}) {
	BeeLogger.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	BeeLogger.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	BeeLogger.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
func Warning(v ...interface{}) {
	BeeLogger.Warning(generateFmtStr(len(v)), v...)
}

// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Warn(v ...interface{}) {
	Warning(v...)
}

func Notice(v ...interface{}) {
	BeeLogger.Notice(generateFmtStr(len(v)), v...)
}

// Info logs a message at info level.
func Informational(v ...interface{}) {
	BeeLogger.Informational(generateFmtStr(len(v)), v...)
}

// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Info(v ...interface{}) {
	Informational(v...)
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	BeeLogger.Debug(generateFmtStr(len(v)), v...)
}

// Trace logs a message at trace level.
// Deprecated: compatibility alias for Warning(), Will be removed in 1.5.0.
func Trace(v ...interface{}) {
	BeeLogger.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
