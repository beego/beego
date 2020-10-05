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
	"io"
	"os"
	"time"

	"github.com/astaxie/beego/pkg/core/governor"
)

var startTime = time.Now()
var pid int

func init() {
	pid = os.Getpid()
}

// ProcessInput parse input command string
func ProcessInput(input string, w io.Writer) {
	governor.ProcessInput(input, w)
}

// MemProf record memory profile in pprof
func MemProf(w io.Writer) {
	governor.MemProf(w)
}

// GetCPUProfile start cpu profile monitor
func GetCPUProfile(w io.Writer) {
	governor.GetCPUProfile(w)
}

// PrintGCSummary print gc information to io.Writer
func PrintGCSummary(w io.Writer) {
	governor.PrintGCSummary(w)
}
