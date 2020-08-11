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
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	beego "github.com/astaxie/beego/pkg"
	"github.com/astaxie/beego/pkg/context"
)

// FilterChainBuilder is an extension point,
// when we want to support some configuration,
// please use this structure
type FilterChainBuilder struct {
}

// FilterChain returns a FilterFunc. The filter will records some metrics
func (builder *FilterChainBuilder) FilterChain(next beego.FilterFunc) beego.FilterFunc {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "beego",
		Subsystem: "http_request",
		ConstLabels: map[string]string{
			"server":  beego.BConfig.ServerName,
			"env":     beego.BConfig.RunMode,
			"appname": beego.BConfig.AppName,
		},
		Help: "The statics info for http request",
	}, []string{"pattern", "method", "status", "duration"})

	prometheus.MustRegister(summaryVec)

	registerBuildInfo()

	return func(ctx *context.Context) {
		startTime := time.Now()
		next(ctx)
		endTime := time.Now()
		go report(endTime.Sub(startTime), ctx, summaryVec)
	}
}

func registerBuildInfo() {
	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "beego",
		Subsystem: "build_info",
		Help:      "The building information",
		ConstLabels: map[string]string{
			"appname":        beego.BConfig.AppName,
			"build_version":  beego.BuildVersion,
			"build_revision": beego.BuildGitRevision,
			"build_status":   beego.BuildStatus,
			"build_tag":      beego.BuildTag,
			"build_time":     strings.Replace(beego.BuildTime, "--", " ", 1),
			"go_version":     beego.GoVersion,
			"git_branch":     beego.GitBranch,
			"start_time":     time.Now().Format("2006-01-02 15:04:05"),
		},
	}, []string{})

	prometheus.MustRegister(buildInfo)
	buildInfo.WithLabelValues().Set(1)
}

func report(dur time.Duration, ctx *context.Context, vec *prometheus.SummaryVec) {
	status := ctx.Output.Status
	ptn := ctx.Input.GetData("RouterPattern").(string)
	ms := dur / time.Millisecond
	vec.WithLabelValues(ptn, ctx.Input.Method(), strconv.Itoa(status), strconv.Itoa(int(ms))).Observe(float64(ms))
}
