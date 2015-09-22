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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/utils"
)

func serverStaticRouter(ctx *context.Context) {
	if ctx.Input.Method() != "GET" && ctx.Input.Method() != "HEAD" {
		return
	}
	requestPath := filepath.Clean(ctx.Input.Request.URL.Path)

	// special processing : favicon.ico/robots.txt  can be in any static dir
	if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {

		if utils.FileExists("./" + requestPath) {
			http.ServeFile(ctx.ResponseWriter, ctx.Request, "./"+requestPath)
			return
		}

		for _, staticDir := range StaticDir {
			file := path.Join(staticDir, requestPath)
			if utils.FileExists(file) {
				http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
				return
			}
		}

		http.NotFound(ctx.ResponseWriter, ctx.Request)
		return
	}

	for prefix, staticDir := range StaticDir {
		if len(prefix) == 0 {
			continue
		}
		if strings.HasPrefix(requestPath, prefix) {
			if len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
				continue
			}
			filePath := path.Join(staticDir, requestPath[len(prefix):])
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				if RunMode == "dev" {
					Warn("Can't find the file:", filePath, err)
				}
				http.NotFound(ctx.ResponseWriter, ctx.Request)
				return
			}
			//if the request is dir and DirectoryIndex is false then
			if fileInfo.IsDir() {
				if !DirectoryIndex {
					exception("403", ctx)
					return
				}
				if ctx.Input.Request.URL.Path[len(ctx.Input.Request.URL.Path)-1] != '/' {
					http.Redirect(ctx.ResponseWriter, ctx.Request, ctx.Input.Request.URL.Path+"/", 302)
					return
				}
			}

			if strings.HasSuffix(requestPath, "/index.html") {
				if utils.FileExists(filePath) {
					fileReader, err := os.Open(filePath)
					if err != nil {
						if RunMode == "dev" {
							Warn("Can't open the file:", filePath, err)
						}
						http.NotFound(ctx.ResponseWriter, ctx.Request)
						return
					}
					defer fileReader.Close()
					http.ServeContent(ctx.ResponseWriter, ctx.Request, filePath, fileInfo.ModTime(), fileReader)
					return
				}
			}

			isStaticFileToCompress := false
			lowerFileName := strings.ToLower(filePath)
			for _, statExtension := range StaticExtensionsToGzip {
				if strings.HasSuffix(lowerFileName, statExtension) {
					isStaticFileToCompress = true
					break
				}
			}

			if !isStaticFileToCompress {
				http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
				return
			}

			//to compress file
			var contentEncoding string
			if EnableGzip {
				contentEncoding = getAcceptEncodingZip(ctx.Request)
			}

			memZipFile, err := openMemZipFile(filePath, contentEncoding)
			if err != nil {
				if RunMode == "dev" {
					Warn("Can't compress the file:", filePath, err)
				}
				http.NotFound(ctx.ResponseWriter, ctx.Request)
				return
			}

			if contentEncoding == "gzip" {
				ctx.Output.Header("Content-Encoding", "gzip")
			} else if contentEncoding == "deflate" {
				ctx.Output.Header("Content-Encoding", "deflate")
			} else {
				ctx.Output.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
			}

			http.ServeContent(ctx.ResponseWriter, ctx.Request, filePath, fileInfo.ModTime(), memZipFile)
		}
	}
}
