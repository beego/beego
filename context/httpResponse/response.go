package httpResponse

import (
	"strconv"

	"net/http"

	beecontext "github.com/astaxie/beego/context"
)

const (
	//BadRequest indicates http error 400
	BadRequest StatusCode = http.StatusBadRequest

	//NotFound indicates http error 404
	NotFound StatusCode = http.StatusNotFound
)

// StatusCode sets the http response status code
type StatusCode int

func (s StatusCode) Error() string {
	return strconv.Itoa(int(s))
}

// Render sets the http status code
func (s StatusCode) Render(ctx *beecontext.Context) {
	ctx.Output.SetStatus(int(s))
}
