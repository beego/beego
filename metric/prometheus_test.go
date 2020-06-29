// Copyright 2020 astaxie
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

package metric

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/astaxie/beego/context"
)

func TestPrometheusMiddleWare(t *testing.T) {
	middleware := PrometheusMiddleWare(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	writer := &context.Response{}
	request := &http.Request{
		URL: &url.URL{
			Host:    "localhost",
			RawPath: "/a/b/c",
		},
		Method: "POST",
	}
	vec := prometheus.NewSummaryVec(prometheus.SummaryOpts{}, []string{"pattern", "method", "status", "duration"})

	report(time.Second, writer, request, vec)
	middleware.ServeHTTP(writer, request)
}
