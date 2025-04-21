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
// Usage can see test cases : TestWrapperFromJsonExample
func WrapperFromJson[T any](
	biz bizFunc[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.BindJSON(&params)
		return
	})
}

// WrapperFromForm  for handling form data in request.
// It binds the form data to the specified type T
// Usage can see test cases : TestWrapperFromFormExample
func WrapperFromForm[T any](
	biz bizFunc[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.BindForm(&params)
		return
	})
}

// Wrapper is use by beego ctx.Bind(any) api
// It binds the data to the specified type T
// Usage can see test cases : TestWrapperExample
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
