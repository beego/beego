// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie
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
	requestPath := path.Clean(ctx.Input.Request.URL.Path)
	for prefix, staticDir := range StaticDir {
		if len(prefix) == 0 {
			continue
		}
		if requestPath == "/favicon.ico" {
			file := path.Join(staticDir, requestPath)
			if utils.FileExists(file) {
				http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
				return
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
			if finfo.IsDir() && !DirectoryIndex {
				middleware.Exception("403", ctx.ResponseWriter, ctx.Request, "403 Forbidden")
				return
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
