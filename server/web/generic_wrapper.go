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
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// bizFunc is a function that wraps a business logic function with a context and parameter extraction function.
type bizFunc[T any] func(ctx *context.Context, param T) (any, error)

// extractFunc is a function that extracts parameters from the context.
type extractFunc[T any] func(ctx *context.Context) (params T, err error)

// WrapperFromJson  for handling JSON in request's body.
// It binds the JSON request body to the specified type T
// Usage can see test cases : ExampleWrapperFromJson
func WrapperFromJson[T any](
	biz bizFunc[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.BindJSON(&params)
		return
	})
}

// WrapperFromForm  for handling form data in request.
// It binds the form data to the specified type T
// Usage can see test cases : ExampleWrapperFromForm
func WrapperFromForm[T any](
	biz bizFunc[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.BindForm(&params)
		return
	})
}

// Wrapper is use by beego ctx.Bind(any) api
// It binds the data to the specified type T
// Usage can see test cases: ExampleWrapper
func Wrapper[T any](
	biz bizFunc[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.Bind(&params)
		return
	})
}

func internalWrapper[T any](
	biz bizFunc[T],
	ef extractFunc[T]) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		params, err := ef(ctx)
		if err != nil {
			logs.Error("err {%v} happen in subject ctx ", err)
			ctx.Abort(400, err.Error())
			return
		}
		res, err := biz(ctx, params)
		if err != nil {
			logs.Error("err {%v} happen in biz ", err)
			ctx.Abort(500, err.Error())
			return
		}
		err = ctx.Resp(res)
		if err != nil {
			logs.Error("err {%v} happen in write response ", err)
			panic(err)
		}
	}
}
