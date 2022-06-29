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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type countTask struct {
	cnt     int
	mockErr error
}

func (c *countTask) GetTask(ctx context.Context) *Task {
	return &Task{}
}

func (c *countTask) GetSpec(ctx context.Context) string {
	return "AAA"
}

func (c *countTask) GetStatus(ctx context.Context) string {
	return "SUCCESS"
}

func (c *countTask) Run(ctx context.Context) error {
	c.cnt++
	return c.mockErr
}

func (c *countTask) SetNext(ctx context.Context, time time.Time) {
}

func (c *countTask) GetNext(ctx context.Context) time.Time {
	return time.Now()
}

func (c *countTask) SetPrev(ctx context.Context, time time.Time) {
}

func (c *countTask) GetPrev(ctx context.Context) time.Time {
	return time.Now()
}

func (c *countTask) GetTimeout(ctx context.Context) time.Duration {
	return 0
}

func TestRunTaskCommand_Execute(t *testing.T) {
	task := &countTask{}
	AddTask("count", task)

	cmd := &runTaskCommand{}

	res := cmd.Execute()
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "task name not passed", res.Error.Error())

	res = cmd.Execute(10)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "parameter is invalid", res.Error.Error())

	res = cmd.Execute("CCCC")
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "task with name CCCC not found", res.Error.Error())

	res = cmd.Execute("count")
	assert.NotNil(t, res)
	assert.True(t, res.IsSuccess())

	task.mockErr = errors.New("mock error")
	res = cmd.Execute("count")
	assert.NotNil(t, res)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "mock error", res.Error.Error())
}

func TestListTaskCommand_Execute(t *testing.T) {
	task := &countTask{}

	cmd := &listTaskCommand{}

	res := cmd.Execute()

	assert.True(t, res.IsSuccess())

	_, ok := res.Content.([][]string)
	assert.True(t, ok)

	AddTask("count", task)

	res = cmd.Execute()

	assert.True(t, res.IsSuccess())

	rl, ok := res.Content.([][]string)
	assert.True(t, ok)
	assert.Equal(t, 1, len(rl))
}
