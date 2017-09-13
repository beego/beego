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
	"bytes"
	"encoding/json"
	"time"
	"fmt"
	"github.com/mcuadros/go-version"
	"runtime"
)

const (
	ApacheFormatPattern = "%s - - [%s] \"%s %d %d\" %f %s %s\n"
	ApacheFormat        = "APACHE_FORMAT"
	JsonFormat          = "JSON_FORMAT"
)

type AccessLogRecord struct {
	RemoteAddr     string        `json:"remote_addr"`
	RequestTime    time.Time     `json:"request_time"`
	RequestMethod  string        `json:"request_method"`
	Request        string        `json:"request"`
	ServerProtocol string        `json:"server_protocol"`
	Host           string        `json:"host"`
	Status         int           `json:"status"`
	BodyBytesSent  int64         `json:"body_bytes_sent"`
	ElapsedTime    time.Duration `json:"elapsed_time"`
	HttpReferrer   string        `json:"http_referrer"`
	HttpUserAgent  string        `json:"http_user_agent"`
	RemoteUser     string        `json:"remote_user"`
}

func (r *AccessLogRecord) JSON() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	currentGoVersion := version.Normalize(runtime.Version()[2:])
	if version.Compare("1.7", currentGoVersion, "<") {
		encoder.SetEscapeHTML(false)
	}
	err := encoder.Encode(r)
	return buffer.Bytes(), err
}

func AccessLog(r *AccessLogRecord, format string) {
	var msg string

	if format == ApacheFormat {
		timeFormatted := r.RequestTime.Format("02/Jan/2006 03:04:05")
		msg = fmt.Sprintf(ApacheFormatPattern, r.RemoteAddr, timeFormatted, r.Request, r.Status, r.BodyBytesSent,
			r.ElapsedTime.Seconds(), r.HttpReferrer, r.HttpUserAgent)
	} else {
		jsonData, err := r.JSON()
		if err != nil {
			msg = fmt.Sprintf(`{"Error": "%s"}`, err)
		} else {
			msg = string(jsonData)
		}
	}
	beeLogger.Debug(msg)
}
