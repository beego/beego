package webx

import (
	"github.com/beego/beego/v2/server/web/context"
)

// WrapperFromJson is a internalWrapper function for handling JSON in request's body.
// usage:
//
//	web.Post("/hello", WrapperFromJson(func(ctx *context.Context, param T) (any, error) {
//		 return param, nil
//	}))
//
// It binds the JSON request body to the specified type T
// See test cases for details
func WrapperFromJson[T any](
	biz func(ctx *context.Context, param T) (any, error)) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context, params *T) error {
		return ctx.BindJSON(params)
	})
}

// WrapperFromForm is a internalWrapper function for handling form data in request.
// usage:
//
//	web.Post("/hello", WrapperFromForm(func(ctx *context.Context, param T) (any, error) {
//		 return param, nil
//	}))
//
// It binds the form data to the specified type T
// See test cases for details
func WrapperFromForm[T any](
	biz func(ctx *context.Context, param T) (any, error)) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context, params *T) error {
		return ctx.BindForm(params)
	})
}

// Wrapper is use by beego ctx.Bind(any) api
// usage:
//
//	web.Post("/hello", Wrapper(func(ctx *context.Context, param T) (any, error) {
//		 return param, nil
//	}))
//
// It binds the data to the specified type T
// See test cases for details
func Wrapper[T any](
	biz func(ctx *context.Context, param T) (any, error)) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context, params *T) error {
		return ctx.Bind(params)
	})
}

// internalWrapper is a core helper function that wraps the business logic and the binding logic.
func internalWrapper[T any](
	biz func(ctx *context.Context, param T) (any, error),
	wf func(ctx *context.Context, params *T) error) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		var params T
		err := wf(ctx, &params)
		if err != nil {
			ctx.Abort(400, err.Error())
			return
		}
		res, err := biz(ctx, params)
		if err != nil {
			ctx.Abort(500, err.Error())
			return
		}
		err = ctx.Resp(res)
		if err != nil {
			panic(err)
		}
	}
}
