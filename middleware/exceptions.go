package middleware

import "fmt"

type HTTPException struct {
	StatusCode  int // http status code 4xx, 5xx
	Description string
}

func (e *HTTPException) Error() string {
	// return `status description`, e.g. `400 Bad Request`
	return fmt.Sprintf("%d %s", e.StatusCode, e.Description)
}

var HTTPExceptionMaps map[int]HTTPException

func init() {
	HTTPExceptionMaps = make(map[int]HTTPException)

	// Normal 4XX HTTP Status
	HTTPExceptionMaps[400] = HTTPException{400, "Bad Request"}
	HTTPExceptionMaps[401] = HTTPException{401, "Unauthorized"}
	HTTPExceptionMaps[403] = HTTPException{403, "Forbidden"}
	HTTPExceptionMaps[404] = HTTPException{404, "Not Found"}
	HTTPExceptionMaps[405] = HTTPException{405, "Method Not Allowed"}

	// Normal 5XX HTTP Status
	HTTPExceptionMaps[500] = HTTPException{500, "Internal Server Error"}
	HTTPExceptionMaps[502] = HTTPException{502, "Bad Gateway"}
	HTTPExceptionMaps[503] = HTTPException{503, "Service Unavailable"}
	HTTPExceptionMaps[504] = HTTPException{504, "Gateway Timeout"}
}
