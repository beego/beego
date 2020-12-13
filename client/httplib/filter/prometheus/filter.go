// Copyright 2020 beego
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

package prometheus

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/beego/beego/client/httplib"
)

type FilterChainBuilder struct {
	summaryVec prometheus.ObserverVec
	AppName    string
	ServerName string
	RunMode    string
}

func (builder *FilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {

	builder.summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "beego",
		Subsystem: "remote_http_request",
		ConstLabels: map[string]string{
			"server":  builder.ServerName,
			"env":     builder.RunMode,
			"appname": builder.AppName,
		},
		Help: "The statics info for remote http requests",
	}, []string{"proto", "scheme", "method", "host", "path", "status", "duration", "isError"})

	return func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {
		startTime := time.Now()
		resp, err := next(ctx, req)
		endTime := time.Now()
		go builder.report(startTime, endTime, ctx, req, resp, err)
		return resp, err
	}
}

func (builder *FilterChainBuilder) report(startTime time.Time, endTime time.Time,
	ctx context.Context, req *httplib.BeegoHTTPRequest, resp *http.Response, err error) {

	proto := req.GetRequest().Proto

	scheme := req.GetRequest().URL.Scheme
	method := req.GetRequest().Method

	host := req.GetRequest().URL.Host
	path := req.GetRequest().URL.Path

	status := -1
	if resp != nil {
		status = resp.StatusCode
	}

	dur := int(endTime.Sub(startTime) / time.Millisecond)

	builder.summaryVec.WithLabelValues(proto, scheme, method, host, path,
		strconv.Itoa(status), strconv.Itoa(dur), strconv.FormatBool(err == nil))
}
