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

package context

import (
	"github.com/beego/beego/v2/server/web/context"
)

// BeegoOutput does work for sending response header.
type BeegoOutput context.BeegoOutput

// NewOutput returns new BeegoOutput.
// it contains nothing now.
func NewOutput() *BeegoOutput {
	return (*BeegoOutput)(context.NewOutput())
}

// Reset init BeegoOutput
func (output *BeegoOutput) Reset(ctx *Context) {
	(*context.BeegoOutput)(output).Reset((*context.Context)(ctx))
}

// Header sets response header item string via given key.
func (output *BeegoOutput) Header(key, val string) {
	(*context.BeegoOutput)(output).Header(key, val)
}

// Body sets response body content.
// if EnableGzip, compress content string.
// it sends out response body directly.
func (output *BeegoOutput) Body(content []byte) error {
	return (*context.BeegoOutput)(output).Body(content)
}

// Cookie sets cookie value via given key.
// others are ordered as cookie's max age time, path,domain, secure and httponly.
func (output *BeegoOutput) Cookie(name string, value string, others ...interface{}) {
	(*context.BeegoOutput)(output).Cookie(name, value, others...)
}

// JSON writes json to response body.
// if encoding is true, it converts utf-8 to \u0000 type.
func (output *BeegoOutput) JSON(data interface{}, hasIndent bool, encoding bool) error {
	return (*context.BeegoOutput)(output).JSON(data, hasIndent, encoding)
}

// YAML writes yaml to response body.
func (output *BeegoOutput) YAML(data interface{}) error {
	return (*context.BeegoOutput)(output).YAML(data)
}

// JSONP writes jsonp to response body.
func (output *BeegoOutput) JSONP(data interface{}, hasIndent bool) error {
	return (*context.BeegoOutput)(output).JSONP(data, hasIndent)
}

// XML writes xml string to response body.
func (output *BeegoOutput) XML(data interface{}, hasIndent bool) error {
	return (*context.BeegoOutput)(output).XML(data, hasIndent)
}

// ServeFormatted serve YAML, XML OR JSON, depending on the value of the Accept header
func (output *BeegoOutput) ServeFormatted(data interface{}, hasIndent bool, hasEncode ...bool) {
	(*context.BeegoOutput)(output).ServeFormatted(data, hasIndent, hasEncode...)
}

// Download forces response for download file.
// it prepares the download response header automatically.
func (output *BeegoOutput) Download(file string, filename ...string) {
	(*context.BeegoOutput)(output).Download(file, filename...)
}

// ContentType sets the content type from ext string.
// MIME type is given in mime package.
func (output *BeegoOutput) ContentType(ext string) {
	(*context.BeegoOutput)(output).ContentType(ext)
}

// SetStatus sets response status code.
// It writes response header directly.
func (output *BeegoOutput) SetStatus(status int) {
	(*context.BeegoOutput)(output).SetStatus(status)
}

// IsCachable returns boolean of this request is cached.
// HTTP 304 means cached.
func (output *BeegoOutput) IsCachable() bool {
	return (*context.BeegoOutput)(output).IsCachable()
}

// IsEmpty returns boolean of this request is empty.
// HTTP 201，204 and 304 means empty.
func (output *BeegoOutput) IsEmpty() bool {
	return (*context.BeegoOutput)(output).IsEmpty()
}

// IsOk returns boolean of this request runs well.
// HTTP 200 means ok.
func (output *BeegoOutput) IsOk() bool {
	return (*context.BeegoOutput)(output).IsOk()
}

// IsSuccessful returns boolean of this request runs successfully.
// HTTP 2xx means ok.
func (output *BeegoOutput) IsSuccessful() bool {
	return (*context.BeegoOutput)(output).IsSuccessful()
}

// IsRedirect returns boolean of this request is redirection header.
// HTTP 301,302,307 means redirection.
func (output *BeegoOutput) IsRedirect() bool {
	return (*context.BeegoOutput)(output).IsRedirect()
}

// IsForbidden returns boolean of this request is forbidden.
// HTTP 403 means forbidden.
func (output *BeegoOutput) IsForbidden() bool {
	return (*context.BeegoOutput)(output).IsForbidden()
}

// IsNotFound returns boolean of this request is not found.
// HTTP 404 means not found.
func (output *BeegoOutput) IsNotFound() bool {
	return (*context.BeegoOutput)(output).IsNotFound()
}

// IsClientError returns boolean of this request client sends error data.
// HTTP 4xx means client error.
func (output *BeegoOutput) IsClientError() bool {
	return (*context.BeegoOutput)(output).IsClientError()
}

// IsServerError returns boolean of this server handler errors.
// HTTP 5xx means server internal error.
func (output *BeegoOutput) IsServerError() bool {
	return (*context.BeegoOutput)(output).IsServerError()
}

// Session sets session item value with given key.
func (output *BeegoOutput) Session(name interface{}, value interface{}) {
	(*context.BeegoOutput)(output).Session(name, value)
}
