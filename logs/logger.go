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

package logs

import (
	"io"
	"sync"
	"time"
)

type logWriter struct {
	sync.Mutex
	writer io.Writer
}

func newLogWriter(wr io.Writer) *logWriter {
	return &logWriter{writer: wr}
}

func (lg *logWriter) println(when time.Time, msg string) {
	lg.Lock()
	h, _ := formatTimeHeader(when)
	lg.writer.Write(append(append(h, msg...), '\n'))
	lg.Unlock()
}

const y1 = `0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999`
const y2 = `0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789`
const mo1 = `000000000111`
const mo2 = `123456789012`
const d1 = `0000000001111111111222222222233`
const d2 = `1234567890123456789012345678901`
const h1 = `000000000011111111112222`
const h2 = `012345678901234567890123`
const mi1 = `000000000011111111112222222222333333333344444444445555555555`
const mi2 = `012345678901234567890123456789012345678901234567890123456789`
const s1 = `000000000011111111112222222222333333333344444444445555555555`
const s2 = `012345678901234567890123456789012345678901234567890123456789`

func formatTimeHeader(when time.Time) ([]byte, int) {
	y, mo, d := when.Date()
	h, mi, s := when.Clock()
	//len("2006/01/02 15:04:05 ")==20
	var buf [20]byte

	//change to '3' after 984 years, LOL
	buf[0] = '2'
	//change to '1' after 84 years,  LOL
	buf[1] = '0'
	buf[2] = y1[y-2000]
	buf[3] = y2[y-2000]
	buf[4] = '/'
	buf[5] = mo1[mo-1]
	buf[6] = mo2[mo-1]
	buf[7] = '/'
	buf[8] = d1[d-1]
	buf[9] = d2[d-1]
	buf[10] = ' '
	buf[11] = h1[h]
	buf[12] = h2[h]
	buf[13] = ':'
	buf[14] = mi1[mi]
	buf[15] = mi2[mi]
	buf[16] = ':'
	buf[17] = s1[s]
	buf[18] = s2[s]
	buf[19] = ' '

	return buf[0:], d
}
