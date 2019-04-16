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

// Logger is the logging interface implementation.
type Logger interface {
	//EnableLog controls whether print the message.
	EnableLog(bool)

	//IsEnabled returns if logger is enabled.
	IsEnabled() bool

	//Print formats using the default formats for its operands and logs the message.
	Print(...interface{})

	//Printf formats according to a format specifier and logs the message.
	Printf(string, ...interface{})
}
