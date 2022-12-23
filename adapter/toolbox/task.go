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

package toolbox

import (
	"context"
	"sort"
	"time"

	"github.com/beego/beego/v2/task"
)

// The bounds for each field.
var (
	AdminTaskList map[string]Tasker
)

const (
	// Set the top bit if a star was included in the expression.
	starBit = 1 << 63
)

// Schedule time taks schedule
type Schedule task.Schedule

// TaskFunc task func type
type TaskFunc func() error

// Tasker task interface
type Tasker interface {
	GetSpec() string
	GetStatus() string
	Run() error
	SetNext(time.Time)
	GetNext() time.Time
	SetPrev(time.Time)
	GetPrev() time.Time
}

// task error
type taskerr struct {
	t       time.Time
	errinfo string
}

// Task task struct
// Deprecated
type Task struct {
	// Deprecated
	Taskname string
	// Deprecated
	Spec *Schedule
	// Deprecated
	SpecStr string
	// Deprecated
	DoFunc TaskFunc
	// Deprecated
	Prev time.Time
	// Deprecated
	Next time.Time
	// Deprecated
	Errlist []*taskerr // like errtime:errinfo
	// Deprecated
	ErrLimit int // max length for the errlist, 0 stand for no limit

	delegate *task.Task
}

// NewTask add new task with name, time and func
func NewTask(tname string, spec string, f TaskFunc) *Task {
	task := task.NewTask(tname, spec, func(ctx context.Context) error {
		return f()
	})
	return &Task{
		delegate: task,
	}
}

// GetSpec get spec string
func (t *Task) GetSpec() string {
	t.initDelegate()

	return t.delegate.GetSpec(context.Background())
}

// GetStatus get current task status
func (t *Task) GetStatus() string {
	t.initDelegate()

	return t.delegate.GetStatus(context.Background())
}

// Run run all tasks
func (t *Task) Run() error {
	t.initDelegate()
	return t.delegate.Run(context.Background())
}

// SetNext set next time for this task
func (t *Task) SetNext(now time.Time) {
	t.initDelegate()
	t.delegate.SetNext(context.Background(), now)
}

// GetNext get the next call time of this task
func (t *Task) GetNext() time.Time {
	t.initDelegate()
	return t.delegate.GetNext(context.Background())
}

// SetPrev set prev time of this task
func (t *Task) SetPrev(now time.Time) {
	t.initDelegate()
	t.delegate.SetPrev(context.Background(), now)
}

// GetPrev get prev time of this task
func (t *Task) GetPrev() time.Time {
	t.initDelegate()
	return t.delegate.GetPrev(context.Background())
}

// six columns mean：
//       second：0-59
//       minute：0-59
//       hour：1-23
//       day：1-31
//       month：1-12
//       week：0-6（0 means Sunday）

// SetCron some signals：
//       *： any time
//       ,：　 separate signal
// 　　    －：duration
//       /n : do as n times of time duration
// ///////////////////////////////////////////////////////
//	0/30 * * * * *                        every 30s
//	0 43 21 * * *                         21:43
//	0 15 05 * * * 　　                     05:15
//	0 0 17 * * *                          17:00
//	0 0 17 * * 1                           17:00 in every Monday
//	0 0,10 17 * * 0,2,3                   17:00 and 17:10 in every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                       17:00 to 17:10 in 1 min duration each time on the first day of month
//	0 0 0 1,15 * 1                        0:00 on the 1st day and 15th day of month
//	0 42 4 1 * * 　 　                     4:42 on the 1st day of month
//	0 0 21 * * 1-6　　                     21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *　           every 10 min duration
//	0 */10 * * * * 　　　　　　              every 10 min duration
//	0 * 1 * * *　　　　　　　　               1:00 to 1:59 in 1 min duration each time
//	0 0 1 * * *　　　　　　　　               1:00
//	0 0 */1 * * *　　　　　　　               0 min of hour in 1 hour duration
//	0 0 * * * *　　　　　　　　               0 min of hour in 1 hour duration
//	0 2 8-20/3 * * *　　　　　　             8:02, 11:02, 14:02, 17:02, 20:02
//	0 30 5 1,15 * *　　　　　　              5:30 on the 1st day and 15th day of month
func (t *Task) SetCron(spec string) {
	t.initDelegate()
	t.delegate.SetCron(spec)
}

func (t *Task) initDelegate() {
	if t.delegate == nil {
		t.delegate = &task.Task{
			Taskname: t.Taskname,
			Spec:     (*task.Schedule)(t.Spec),
			SpecStr:  t.SpecStr,
			DoFunc: func(ctx context.Context) error {
				return t.DoFunc()
			},
			Prev:     t.Prev,
			Next:     t.Next,
			ErrLimit: t.ErrLimit,
		}
	}
}

// Next set schedule to next time
func (s *Schedule) Next(t time.Time) time.Time {
	return (*task.Schedule)(s).Next(t)
}

// StartTask start all tasks
func StartTask() {
	task.StartTask()
}

// StopTask stop all tasks
func StopTask() {
	task.StopTask()
}

// AddTask add task with name
func AddTask(taskname string, t Tasker) {
	task.AddTask(taskname, &oldToNewAdapter{delegate: t})
}

// DeleteTask delete task with name
func DeleteTask(taskname string) {
	task.DeleteTask(taskname)
}

// ClearTask clear all tasks
func ClearTask() {
	task.ClearTask()
}

// MapSorter sort map for tasker
type MapSorter task.MapSorter

// NewMapSorter create new tasker map
func NewMapSorter(m map[string]Tasker) *MapSorter {
	newTaskerMap := make(map[string]task.Tasker, len(m))

	for key, value := range m {
		newTaskerMap[key] = &oldToNewAdapter{
			delegate: value,
		}
	}

	return (*MapSorter)(task.NewMapSorter(newTaskerMap))
}

// Sort sort tasker map
func (ms *MapSorter) Sort() {
	sort.Sort(ms)
}

func (ms *MapSorter) Len() int { return len(ms.Keys) }

func (ms *MapSorter) Less(i, j int) bool {
	if ms.Vals[i].GetNext(context.Background()).IsZero() {
		return false
	}
	if ms.Vals[j].GetNext(context.Background()).IsZero() {
		return true
	}
	return ms.Vals[i].GetNext(context.Background()).Before(ms.Vals[j].GetNext(context.Background()))
}

func (ms *MapSorter) Swap(i, j int) {
	ms.Vals[i], ms.Vals[j] = ms.Vals[j], ms.Vals[i]
	ms.Keys[i], ms.Keys[j] = ms.Keys[j], ms.Keys[i]
}

func init() {
	AdminTaskList = make(map[string]Tasker)
}

type oldToNewAdapter struct {
	delegate Tasker
}

func (o *oldToNewAdapter) GetSpec(ctx context.Context) string {
	return o.delegate.GetSpec()
}

func (o *oldToNewAdapter) GetStatus(ctx context.Context) string {
	return o.delegate.GetStatus()
}

func (o *oldToNewAdapter) Run(ctx context.Context) error {
	return o.delegate.Run()
}

func (o *oldToNewAdapter) SetNext(ctx context.Context, t time.Time) {
	o.delegate.SetNext(t)
}

func (o *oldToNewAdapter) GetNext(ctx context.Context) time.Time {
	return o.delegate.GetNext()
}

func (o *oldToNewAdapter) SetPrev(ctx context.Context, t time.Time) {
	o.delegate.SetPrev(t)
}

func (o *oldToNewAdapter) GetPrev(ctx context.Context) time.Time {
	return o.delegate.GetPrev()
}

func (o *oldToNewAdapter) GetTimeout(ctx context.Context) time.Duration {
	return 0
}
