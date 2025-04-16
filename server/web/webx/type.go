package webx

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

type BizFunc[T any] func(ctx *context.Context, param T) (any, error)
type ExtractFunc[T any] func(ctx *context.Context) (params T, err error)

// WrapperFromJson is a internalWrapper function for handling JSON in request's body.
// It binds the JSON request body to the specified type T
// See test cases for details
func WrapperFromJson[T any](
	biz BizFunc[T],
	opts ...Option[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.BindJSON(&params)
		return
	}, opts...)
}

// WrapperFromForm is a internalWrapper function for handling form data in request.
// It binds the form data to the specified type T
// See test cases for details
func WrapperFromForm[T any](
	biz BizFunc[T],
	opts ...Option[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.BindForm(&params)
		return
	}, opts...)
}

// Wrapper is use by beego ctx.Bind(any) api
// It binds the data to the specified type T
// See test cases for details
func Wrapper[T any](
	biz BizFunc[T],
	opts ...Option[T]) func(ctx *context.Context) {
	return internalWrapper(biz, func(ctx *context.Context) (params T, err error) {
		err = ctx.Bind(&params)
		return
	}, opts...)
}

func internalWrapper[T any](
	biz BizFunc[T],
	sf ExtractFunc[T],
	opts ...Option[T]) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		params, err := sf(ctx)
		if err != nil {
			logs.Error("err {%v} happen in subject ctx ", err)
			ctx.Abort(400, err.Error())
			return
		}

		for _, opt := range opts {
			err = opt(ctx, &params)
			if err != nil {
				logs.Error("err {%v} happen in subject opts ctx ", err)
				ctx.Abort(400, err.Error())
				return
			}
		}

		res, err := biz(ctx, params)
		if err != nil {
			logs.Error("err {%v} happen in biz ", err)
			ctx.Abort(500, err.Error())
			return
		}
		err = ctx.Resp(res)
		if err != nil {
			panic(err)
		}
	}
}

// Option options to T
type Option[T any] func(ctx *context.Context, t *T) error

// InjectSession is a option to inject session into T
func InjectSession[T any]() Option[T] {
	return func(ctx *context.Context, t *T) error {
		if holder, ok := any(t).(SessionHolder); ok {
			store, err := ctx.Session()
			if err != nil {
				return err
			}
			holder.setSession(store)
		}
		return nil
	}
}

type SessionHolder interface {
	GetSession() session.Store
	setSession(s session.Store)
}
