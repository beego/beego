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

package beego

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/middleware"
	"github.com/astaxie/beego/utils"
)

func serverStaticRouter(ctx *context.Context) {
	if ctx.Input.Method() != "GET" && ctx.Input.Method() != "HEAD" {
		return
	}
	requestPath := path.Clean(ctx.Input.Request.URL.Path)
	i := 0
	for prefix, staticDir := range StaticDir {
		if len(prefix) == 0 {
			continue
		}
		if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
			file := path.Join(staticDir, requestPath)
			if utils.FileExists(file) {
				http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
				return
			} else {
				i++
				if i == len(StaticDir) {
					http.NotFound(ctx.ResponseWriter, ctx.Request)
					return
				} else {
					continue
				}
			}
		}
		if strings.HasPrefix(requestPath, prefix) {
			if len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
				continue
			}
			file := path.Join(staticDir, requestPath[len(prefix):])
			finfo, err := os.Stat(file)
			if err != nil {
				if RunMode == "dev" {
					Warn(err)
				}
				http.NotFound(ctx.ResponseWriter, ctx.Request)
				return
			}
			//if the request is dir and DirectoryIndex is false then
			if finfo.IsDir() {
				if !DirectoryIndex {
					middleware.Exception("403", ctx.ResponseWriter, ctx.Request, "403 Forbidden")
					return
				} else if ctx.Input.Request.URL.Path[len(ctx.Input.Request.URL.Path)-1] != '/' {
					http.Redirect(ctx.ResponseWriter, ctx.Request, ctx.Input.Request.URL.Path+"/", 302)
					return
				}
			} else if strings.HasSuffix(requestPath, "/index.html") {
				file := path.Join(staticDir, requestPath)
				if utils.FileExists(file) {
					http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
					return
				}
			}

			//This block obtained from (https://github.com/smithfox/beego) - it should probably get merged into astaxie/beego after a pull request
			isStaticFileToCompress := false
			if StaticExtensionsToGzip != nil && len(StaticExtensionsToGzip) > 0 {
				for _, statExtension := range StaticExtensionsToGzip {
					if strings.HasSuffix(strings.ToLower(file), strings.ToLower(statExtension)) {
						isStaticFileToCompress = true
						break
					}
				}
			}

			if isStaticFileToCompress {
				var contentEncoding string
				if EnableGzip {
					contentEncoding = getAcceptEncodingZip(ctx.Request)
				}

				memzipfile, err := openMemZipFile(file, contentEncoding)
				if err != nil {
					return
				}

				if contentEncoding == "gzip" {
					ctx.Output.Header("Content-Encoding", "gzip")
				} else if contentEncoding == "deflate" {
					ctx.Output.Header("Content-Encoding", "deflate")
				} else {
					ctx.Output.Header("Content-Length", strconv.FormatInt(finfo.Size(), 10))
				}

				http.ServeContent(ctx.ResponseWriter, ctx.Request, file, finfo.ModTime(), memzipfile)

			} else {
				http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
			}
			return
		}
	}
}
