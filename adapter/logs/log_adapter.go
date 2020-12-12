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
	"time"

	"github.com/astaxie/beego/core/logs"
)

type oldToNewAdapter struct {
	old Logger
}

func (o *oldToNewAdapter) Init(config string) error {
	return o.old.Init(config)
}

func (o *oldToNewAdapter) WriteMsg(lm *logs.LogMsg) error {
	return o.old.WriteMsg(lm.When, lm.OldStyleFormat(), lm.Level)
}

func (o *oldToNewAdapter) Destroy() {
	o.old.Destroy()
}

func (o *oldToNewAdapter) Flush() {
	o.old.Flush()
}

func (o *oldToNewAdapter) SetFormatter(f logs.LogFormatter) {
	panic("unsupported operation, you should not invoke this method")
}

type newToOldAdapter struct {
	n logs.Logger
}

func (n *newToOldAdapter) Init(config string) error {
	return n.n.Init(config)
}

func (n *newToOldAdapter) WriteMsg(when time.Time, msg string, level int) error {
	return n.n.WriteMsg(&logs.LogMsg{
		When:  when,
		Msg:   msg,
		Level: level,
	})
}

func (n *newToOldAdapter) Destroy() {
	panic("implement me")
}

func (n *newToOldAdapter) Flush() {
	panic("implement me")
}
