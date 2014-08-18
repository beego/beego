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
