package response

import (
	beecontext "github.com/astaxie/beego/context"
)

func Json(value interface{}, encoding ...bool) Renderer {
	return rendererFunc(func(ctx *beecontext.Context) {
		var (
			hasIndent   = true
			hasEncoding = false
		)
		//TODO: need access to BConfig :(
		// if BConfig.RunMode == PROD {
		// 	hasIndent = false
		// }
		if len(encoding) > 0 && encoding[0] {
			hasEncoding = true
		}
		ctx.Output.JSON(value, hasIndent, hasEncoding)
	})
}

func errorRenderer(err error) Renderer {
	return rendererFunc(func(ctx *beecontext.Context) {
		ctx.Output.SetStatus(500)
		ctx.WriteString(err.Error())
	})
}

func Redirect(localurl string) statusCodeWithRender {
	return statusCodeWithRender{302, func(ctx *beecontext.Context) {
		ctx.Redirect(302, localurl)
	}}
}

func RenderMethodResult(result interface{}, ctx *beecontext.Context) {
	if result != nil {
		renderer, ok := result.(Renderer)
		if !ok {
			err, ok := result.(error)
			if ok {
				renderer = errorRenderer(err)
			} else {
				renderer = Json(result)
			}
		}
		renderer.Render(ctx)
	}
}
