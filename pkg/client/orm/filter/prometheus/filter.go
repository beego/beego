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
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/astaxie/beego/pkg/client/orm"
)

// FilterChainBuilder is an extension point,
// when we want to support some configuration,
// please use this structure
// this Filter's behavior looks a little bit strange
// for example:
// if we want to records the metrics of QuerySetter
// actually we only records metrics of invoking "QueryTable" and "QueryTableWithCtx"
type FilterChainBuilder struct {
	summaryVec prometheus.ObserverVec
	AppName    string
	ServerName string
	RunMode    string
}

func (builder *FilterChainBuilder) FilterChain(next orm.Filter) orm.Filter {

	builder.summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "beego",
		Subsystem: "orm_operation",
		ConstLabels: map[string]string{
			"server":  builder.ServerName,
			"env":     builder.RunMode,
			"appname": builder.AppName,
		},
		Help: "The statics info for orm operation",
	}, []string{"method", "name", "duration", "insideTx", "txName"})

	return func(ctx context.Context, inv *orm.Invocation) []interface{} {
		startTime := time.Now()
		res := next(ctx, inv)
		endTime := time.Now()
		dur := (endTime.Sub(startTime)) / time.Millisecond

		// if the TPS is too large, here may be some problem
		// thinking about using goroutine pool
		go builder.report(ctx, inv, dur)
		return res
	}
}

func (builder *FilterChainBuilder) report(ctx context.Context, inv *orm.Invocation, dur time.Duration) {
	// start a transaction, we don't record it
	if strings.HasPrefix(inv.Method, "Begin") {
		return
	}
	if inv.Method == "Commit" || inv.Method == "Rollback" {
		builder.reportTxn(ctx, inv)
		return
	}
	builder.summaryVec.WithLabelValues(inv.Method, inv.GetTableName(), strconv.Itoa(int(dur)),
		strconv.FormatBool(inv.InsideTx), inv.TxName)
}

func (builder *FilterChainBuilder) reportTxn(ctx context.Context, inv *orm.Invocation) {
	dur := time.Now().Sub(inv.TxStartTime) / time.Millisecond
	builder.summaryVec.WithLabelValues(inv.Method, inv.TxName, strconv.Itoa(int(dur)),
		strconv.FormatBool(inv.InsideTx), inv.TxName)
}
