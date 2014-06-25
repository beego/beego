// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package middleware

import "fmt"

// http exceptions
type HTTPException struct {
	StatusCode  int // http status code 4xx, 5xx
	Description string
}

// return http exception error string, e.g. "400 Bad Request".
func (e *HTTPException) Error() string {
	return fmt.Sprintf("%d %s", e.StatusCode, e.Description)
}

// map of http exceptions for each http status code int.
// defined 400,401,403,404,405,500,502,503 and 504 default.
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
