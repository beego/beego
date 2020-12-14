package context

import (
	"net/http"
	"strconv"
)

const (
	// BadRequest indicates HTTP error 400
	BadRequest StatusCode = http.StatusBadRequest

	// NotFound indicates HTTP error 404
	NotFound StatusCode = http.StatusNotFound
)

// StatusCode sets the HTTP response status code
type StatusCode int

func (s StatusCode) Error() string {
	return strconv.Itoa(int(s))
}

// Render sets the HTTP status code
func (s StatusCode) Render(ctx *Context) {
	ctx.Output.SetStatus(int(s))
}
