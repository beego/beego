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

package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBeeLoggerDelLogger(t *testing.T) {
	prefix := "My-Cus"
	l := GetLogger(prefix)
	assert.NotNil(t, l)
	l.Print("hello")

	GetLogger().Print("hello")
	SetPrefix("aaa")
	Info("hello")
}

type mockLogger struct {
	*logWriter
	WriteCost time.Duration `json:"write_cost"` // Simulated log writing time consuming
	writeCnt  int           // Count add 1 when writing log success, just for test result
}

func NewMockLogger() Logger {
	return &mockLogger{
		//logWriter: newLogWriter(ansicolor.NewAnsiColorWriter(os.Stdout)),
		logWriter: &logWriter{writer: io.Discard},
	}
}

func (m *mockLogger) Init(config string) error {
	return json.Unmarshal([]byte(config), m)
}

func (m *mockLogger) WriteMsg(lm *LogMsg) error {
	m.Lock()
	msg := lm.Msg
	msg += "\n"

	time.Sleep(m.WriteCost)
	if _, err := m.writer.Write([]byte(msg)); err != nil {
		return err
	}

	m.writeCnt++
	m.Unlock()
	return nil
}

func (m *mockLogger) GetCnt() int {
	return m.writeCnt
}

func (m *mockLogger) Destroy()                    {}
func (m *mockLogger) Flush()                      {}
func (m *mockLogger) SetFormatter(f LogFormatter) {}

func TestBeeLogger_AsyncNonBlockWrite(t *testing.T) {
	testCases := []struct {
		name         string
		before       func()
		after        func()
		msgLen       int64
		writeCost    time.Duration
		sendInterval time.Duration
		writeCnt     int
	}{
		{
			// Write log time is less than send log time, no blocking
			name: "mock1",
			after: func() {
				_ = beeLogger.DelLogger("mock1")
			},
			msgLen:       5,
			writeCnt:     10,
			writeCost:    200 * time.Millisecond,
			sendInterval: 300 * time.Millisecond,
		},
		{
			// Write log time is less than send log time, discarded when blocking
			name: "mock2",
			after: func() {
				_ = beeLogger.DelLogger("mock2")
			},
			writeCnt:     5,
			msgLen:       5,
			writeCost:    200 * time.Millisecond,
			sendInterval: 10 * time.Millisecond,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Register(tc.name, NewMockLogger)
			err := beeLogger.SetLogger(tc.name, fmt.Sprintf(`{"write_cost": %d}`, tc.writeCost))
			assert.Nil(t, err)

			l := beeLogger
			l = beeLogger.Async(tc.msgLen)
			l.AsyncNonBlockWrite()

			for i := 0; i < 10; i++ {
				time.Sleep(tc.sendInterval)
				l.Info(fmt.Sprintf("----%d----", i))
			}
			time.Sleep(1 * time.Second)
			assert.Equal(t, tc.writeCnt, l.outputs[0].Logger.(*mockLogger).writeCnt)
			tc.after()
		})
	}
}
