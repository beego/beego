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

package admin

import (
	"errors"
)

// Command is an experimental interface
// We try to use this to decouple modules
// All other modules depends on this, and they register the command they support
// We may change the API in the future, so be careful about this.
type Command interface {
	Execute(params ...interface{}) *Result
}

var CommandNotFound = errors.New("Command not found")

type Result struct {
	// Status is the same as http.Status
	Status  int
	Error   error
	Content interface{}
}

func (r *Result) IsSuccess() bool {
	return r.Status >= 200 && r.Status < 300
}

// CommandRegistry stores all commands
// name => command
type moduleCommands map[string]Command

// Get returns command with the name
func (m moduleCommands) Get(name string) Command {
	c, ok := m[name]
	if ok {
		return c
	}
	return &doNothingCommand{}
}

// module name => moduleCommand
type commandRegistry map[string]moduleCommands

// Get returns module's commands
func (c commandRegistry) Get(moduleName string) moduleCommands {
	if mcs, ok := c[moduleName]; ok {
		return mcs
	}
	res := make(moduleCommands)
	c[moduleName] = res
	return res
}

var cmdRegistry = make(commandRegistry)

// RegisterCommand is not thread-safe
// do not use it in concurrent case
func RegisterCommand(module string, commandName string, command Command) {
	cmdRegistry.Get(module)[commandName] = command
}

func GetCommand(module string, cmdName string) Command {
	return cmdRegistry.Get(module).Get(cmdName)
}

type doNothingCommand struct{}

func (d *doNothingCommand) Execute(params ...interface{}) *Result {
	return &Result{
		Status: 404,
		Error:  CommandNotFound,
	}
}
