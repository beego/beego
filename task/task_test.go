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

package task

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	tk := NewTask("taska", "0/30 * * * * *", func(ctx context.Context) error {
		fmt.Println("hello world")
		return nil
	})
	err := tk.Run(nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0/30 * * * * *", tk.GetSpec(context.Background()))
	m.AddTask("taska", tk)
	m.StartTask()
	time.Sleep(3 * time.Second)
	assert.True(t, len(tk.GetStatus(context.Background())) == 0)
	m.StopTask()
}

func TestModifyTaskListAfterRunning(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	tk := NewTask("taskb", "0/30 * * * * *", func(ctx context.Context) error {
		fmt.Println("hello world")
		return nil
	})
	err := tk.Run(nil)
	if err != nil {
		t.Fatal(err)
	}
	m.AddTask("taskb", tk)
	m.StartTask()
	go func() {
		m.DeleteTask("taskb")
	}()
	go func() {
		m.AddTask("taskb1", tk)
	}()

	time.Sleep(3 * time.Second)
	m.StopTask()
}

func TestSpec(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tk1 := NewTask("tk1", "0 12 * * * *", func(ctx context.Context) error { fmt.Println("tk1"); return nil })
	tk2 := NewTask("tk2", "0,10,20 * * * * *", func(ctx context.Context) error { fmt.Println("tk2"); wg.Done(); return nil })
	tk3 := NewTask("tk3", "0 10 * * * *", func(ctx context.Context) error { fmt.Println("tk3"); wg.Done(); return nil })

	m.AddTask("tk1", tk1)
	m.AddTask("tk2", tk2)
	m.AddTask("tk3", tk3)
	m.StartTask()
	defer m.StopTask()

	select {
	case <-time.After(200 * time.Second):
		t.FailNow()
	case <-wait(wg):
	}
}

func TestTimeout(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	once1, once2 := sync.Once{}, sync.Once{}

	tk1 := NewTask("tk1", "0/10 * * ? * *",
		func(ctx context.Context) error {
			time.Sleep(4 * time.Second)
			select {
			case <-ctx.Done():
				once1.Do(func() {
					fmt.Println("tk1 done")
					wg.Done()
				})
				return errors.New("timeout")
			default:
			}
			return nil
		}, TimeoutOption(3*time.Second),
	)

	tk2 := NewTask("tk2", "0/11 * * ? * *",
		func(ctx context.Context) error {
			time.Sleep(4 * time.Second)
			select {
			case <-ctx.Done():
				return errors.New("timeout")
			default:
				once2.Do(func() {
					fmt.Println("tk2 done")
					wg.Done()
				})
			}
			return nil
		},
	)

	m.AddTask("tk1", tk1)
	m.AddTask("tk2", tk2)
	m.StartTask()
	defer m.StopTask()

	select {
	case <-time.After(19 * time.Second):
		t.Error("TestTimeout failed")
	case <-wait(wg):
	}
}

func TestTask_Run(t *testing.T) {
	cnt := -1
	task := func(ctx context.Context) error {
		cnt++
		fmt.Printf("Hello, world! %d \n", cnt)
		return errors.New(fmt.Sprintf("Hello, world! %d", cnt))
	}
	tk := NewTask("taska", "0/30 * * * * *", task)
	for i := 0; i < 200; i++ {
		e := tk.Run(nil)
		assert.NotNil(t, e)
	}

	l := tk.Errlist
	assert.Equal(t, 100, len(l))
	assert.Equal(t, "Hello, world! 100", l[0].errinfo)
	assert.Equal(t, "Hello, world! 101", l[1].errinfo)
}

func TestCrudTask(t *testing.T) {
	m := newTaskManager()
	m.AddTask("my-task1", NewTask("my-task1", "0/30 * * * * *", func(ctx context.Context) error {
		return nil
	}))

	m.AddTask("my-task2", NewTask("my-task2", "0/30 * * * * *", func(ctx context.Context) error {
		return nil
	}))

	m.DeleteTask("my-task1")
	assert.Equal(t, 1, len(m.adminTaskList))

	m.ClearTask()
	assert.Equal(t, 0, len(m.adminTaskList))
}

func TestGracefulShutdown(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()
	waitDone := atomic.Value{}
	waitDone.Store(false)
	tk := NewTask("everySecond", "* * * * * *", func(ctx context.Context) error {
		fmt.Println("hello world")
		time.Sleep(2 * time.Second)
		waitDone.Store(true)
		return nil
	})
	m.AddTask("taska", tk)
	m.StartTask()
	time.Sleep(1 * time.Second)
	shutdown := m.GracefulShutdown()
	assert.False(t, waitDone.Load().(bool))
	<-shutdown
	assert.True(t, waitDone.Load().(bool))
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func TestGetAllTasks(t *testing.T) {
	m := newTaskManager()
	defer m.ClearTask()

	tk := NewTask("task1", "0/30 * * * * *", func(ctx context.Context) error {
		return nil
	})

	tk2 := NewTask("task2", "0/40 * * * * *", func(ctx context.Context) error {
		return nil
	})

	m.AddTask("task1", tk)
	m.AddTask("task2", tk2)

	tasks := m.GetAllTasks()
	total := len(tasks)

	assert.Equal(t, 2, total)
}
