// Copyright 2018 The casbin Authors. All Rights Reserved.
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

package log

import "log"

// DefaultLogger is the implementation for a Logger using golang log.
type DefaultLogger struct {
	enable bool
}

func (l *DefaultLogger) EnableLog(enable bool) {
	l.enable = enable
}

func (l *DefaultLogger) IsEnabled() bool {
	return l.enable
}

func (l *DefaultLogger) Print(v ...interface{}) {
	if l.enable {
		log.Print(v...)
	}
}

func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	if l.enable {
		log.Printf(format, v...)
	}
}
