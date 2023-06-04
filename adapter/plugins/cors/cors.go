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

// Package cors provides handlers to enable CORS support.
// Usage
//
//	import (
//		"github.com/beego/beego/v2"
//		"github.com/beego/beego/v2/server/web/filter/cors"
//
// )
//
//	func main() {
//		// CORS for https://foo.* origins, allowing:
//		// - PUT and PATCH methods
//		// - Origin header
//		// - Credentials share
//		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
//			AllowOrigins:     []string{"https://*.foo.com"},
//			AllowMethods:     []string{"PUT", "PATCH"},
//			AllowHeaders:     []string{"Origin"},
//			ExposeHeaders:    []string{"Content-Length"},
//			AllowCredentials: true,
//		}))
//		beego.Run()
//	}
package cors

import (
	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/context"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

// Options represents Access Control options.
type Options cors.Options

// Header converts options into CORS headers.
func (o *Options) Header(origin string) (headers map[string]string) {
	return (*cors.Options)(o).Header(origin)
}

// PreflightHeader converts options into CORS headers for a preflight response.
func (o *Options) PreflightHeader(origin, rMethod, rHeaders string) (headers map[string]string) {
	return (*cors.Options)(o).PreflightHeader(origin, rMethod, rHeaders)
}

// IsOriginAllowed looks up if the origin matches one of the patterns
// generated from Options.AllowOrigins patterns.
func (o *Options) IsOriginAllowed(origin string) bool {
	return (*cors.Options)(o).IsOriginAllowed(origin)
}

// Allow enables CORS for requests those match the provided options.
func Allow(opts *Options) beego.FilterFunc {
	f := cors.Allow((*cors.Options)(opts))
	return func(c *context.Context) {
		f((*beecontext.Context)(c))
	}
}
