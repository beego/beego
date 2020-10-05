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

package task

import (
	"context"
	"fmt"
	"html/template"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/pkg/infrastructure/governor"
)

type listTaskCommand struct {
}

func (l *listTaskCommand) Execute(params ...interface{}) *governor.Result {
	resultList := make([][]string, 0, len(globalTaskManager.adminTaskList))
	for tname, tk := range globalTaskManager.adminTaskList {
		result := []string{
			template.HTMLEscapeString(tname),
			template.HTMLEscapeString(tk.GetSpec(nil)),
			template.HTMLEscapeString(tk.GetStatus(nil)),
			template.HTMLEscapeString(tk.GetPrev(context.Background()).String()),
		}
		resultList = append(resultList, result)
	}

	return &governor.Result{
		Status:  200,
		Content: resultList,
	}
}

type runTaskCommand struct {
}

func (r *runTaskCommand) Execute(params ...interface{}) *governor.Result {
	if len(params) == 0 {
		return &governor.Result{
			Status: 400,
			Error:  errors.New("task name not passed"),
		}
	}

	tn, ok := params[0].(string)

	if !ok {
		return &governor.Result{
			Status: 400,
			Error:  errors.New("parameter is invalid"),
		}
	}

	if t, ok := globalTaskManager.adminTaskList[tn]; ok {
		err := t.Run(context.Background())
		if err != nil {
			return &governor.Result{
				Status: 500,
				Error:  err,
			}
		}
		return &governor.Result{
			Status:  200,
			Content: t.GetStatus(context.Background()),
		}
	} else {
		return &governor.Result{
			Status: 400,
			Error:  errors.New(fmt.Sprintf("task with name %s not found", tn)),
		}
	}

}

func registerCommands() {
	governor.RegisterCommand("task", "list", &listTaskCommand{})
	governor.RegisterCommand("task", "run", &runTaskCommand{})
}
