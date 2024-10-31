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

package web

import (
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

// FilterChain is different from pure FilterFunc
// when you use this, you must invoke next(ctx) inside the FilterFunc which is returned
// And all those FilterChain will be invoked before other FilterFunc
type FilterChain func(next FilterFunc) FilterFunc

// FilterFunc defines a filter function which is invoked before the controller handler is executed.
// It's a alias of HandleFunc
// In fact, the HandleFunc is the last Filter. This is the truth
type FilterFunc = HandleFunc

// FilterRouter defines a filter operation which is invoked before the controller handler is executed.
// It can match the URL against a pattern, and execute a filter function
// when a request with a matching URL arrives.
type FilterRouter struct {
	filterFunc     FilterFunc
	next           *FilterRouter
	tree           *Tree
	pattern        string
	returnOnOutput bool
	resetParams    bool
}

// params is for:
//  1. setting the returnOnOutput value (false allows multiple filters to execute)
//  2. determining whether or not params need to be reset.
func newFilterRouter(pattern string, filter FilterFunc, opts ...FilterOpt) *FilterRouter {
	mr := &FilterRouter{
		tree:       NewTree(),
		pattern:    pattern,
		filterFunc: filter,
	}

	fos := &filterOpts{
		returnOnOutput: true,
	}

	for _, o := range opts {
		o(fos)
	}

	if !fos.routerCaseSensitive {
		mr.pattern = strings.ToLower(pattern)
	}

	mr.returnOnOutput = fos.returnOnOutput
	mr.resetParams = fos.resetParams
	mr.tree.AddRouter(pattern, true)
	return mr
}

// filter will check whether we need to execute the filter logic
// return (started, done)
func (f *FilterRouter) filter(ctx *context.Context, urlPath string, preFilterParams map[string]string) (bool, bool) {
	if f.returnOnOutput && ctx.ResponseWriter.Started {
		return true, true
	}
	if f.resetParams {
		preFilterParams = ctx.Input.Params()
	}
	if ok := f.ValidRouter(urlPath, ctx); ok {
		f.filterFunc(ctx)
		if f.resetParams {
			ctx.Input.ResetParams()
			for k, v := range preFilterParams {
				ctx.Input.SetParam(k, v)
			}
		}
	} else if f.next != nil {
		return f.next.filter(ctx, urlPath, preFilterParams)
	}
	if f.returnOnOutput && ctx.ResponseWriter.Started {
		return true, true
	}
	return false, false
}

// ValidRouter checks if the current request is matched by this filter.
// If the request is matched, the values of the URL parameters defined
// by the filter pattern are also returned.
func (f *FilterRouter) ValidRouter(url string, ctx *context.Context) bool {
	isOk := f.tree.Match(url, ctx)
	if isOk != nil {
		if b, ok := isOk.(bool); ok {
			return b
		}
	}
	return false
}

type filterOpts struct {
	returnOnOutput      bool
	resetParams         bool
	routerCaseSensitive bool
}

type FilterOpt func(opts *filterOpts)

func WithReturnOnOutput(ret bool) FilterOpt {
	return func(opts *filterOpts) {
		opts.returnOnOutput = ret
	}
}

func WithResetParams(reset bool) FilterOpt {
	return func(opts *filterOpts) {
		opts.resetParams = reset
	}
}

func WithCaseSensitive(sensitive bool) FilterOpt {
	return func(opts *filterOpts) {
		opts.routerCaseSensitive = sensitive
	}
}
