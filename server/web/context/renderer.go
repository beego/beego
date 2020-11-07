package context

// Renderer defines a http response renderer
type Renderer interface {
	Render(ctx *Context)
}

type rendererFunc func(ctx *Context)

func (f rendererFunc) Render(ctx *Context) {
	f(ctx)
}
