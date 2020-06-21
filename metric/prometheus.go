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
	"reflect"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func PrometheusMiddleWare(next http.Handler) http.Handler {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      beego.BConfig.AppName,
		Subsystem: "http_request",
		ConstLabels: map[string]string{
			"server": beego.BConfig.ServerName,
			"env":    beego.BConfig.RunMode,
		},
		Help: "The statics info for http request",
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(summaryVec)

	return http.HandlerFunc(func(writer http.ResponseWriter, q *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, q)
		end := time.Now()
		go report(end.Sub(start), writer, q, summaryVec)
	})
}

func report(dur time.Duration, writer http.ResponseWriter, q *http.Request, vec *prometheus.SummaryVec) {
	ctrl := beego.BeeApp.Handlers
	ctx := ctrl.GetContext()
	ctx.Reset(writer, q)
	defer ctrl.GiveBackContext(ctx)

	// We cannot read the status code from q.Response.StatusCode
	// since the http server does not set q.Response. So q.Response is nil
	// Thus, we use reflection to read the status from writer whose concrete type is http.response
	responseVal := reflect.ValueOf(writer).Elem()
	field := responseVal.FieldByName("status")
	status := -1
	if field.IsValid() && field.Kind() == reflect.Int {
		status = int(field.Int())
	}
	ptn := "UNKNOWN"
	if rt, found := ctrl.FindRouter(ctx); found {
		ptn = rt.GetPattern()
	} else {
		logs.Warn("we can not find the router info for this request, so request will be recorded as UNKNOWN: " + q.URL.String())
	}
	vec.WithLabelValues(ptn, q.Method, strconv.Itoa(status)).Observe(float64(dur / time.Millisecond))
}
