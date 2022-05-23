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

package adapter

import (
	"html/template"
	"net/url"
	"time"

	"github.com/beego/beego/v2/server/web"
)

const (
	formatTime      = "15:04:05"
	formatDate      = "2006-01-02"
	formatDateTime  = "2006-01-02 15:04:05"
	formatDateTimeT = "2006-01-02T15:04:05"
)

// Substr returns the substr from start to length.
func Substr(s string, start, length int) string {
	return web.Substr(s, start, length)
}

// HTML2str returns escaping text convert from html.
func HTML2str(html string) string {
	return web.HTML2str(html)
}

// DateFormat takes a time and a layout string and returns a string with the formatted date. Used by the template parser as "dateformat"
func DateFormat(t time.Time, layout string) (datestring string) {
	return web.DateFormat(t, layout)
}

// DateParse Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	return web.DateParse(dateString, format)
}

// Date takes a PHP like date func to Go's time format.
func Date(t time.Time, format string) string {
	return web.Date(t, format)
}

// Compare is a quick and dirty comparison function. It will convert whatever you give it to strings and see if the two values are equal.
// Whitespace is trimmed. Used by the template parser as "eq".
func Compare(a, b interface{}) (equal bool) {
	return web.Compare(a, b)
}

// CompareNot !Compare
func CompareNot(a, b interface{}) (equal bool) {
	return web.CompareNot(a, b)
}

// NotNil the same as CompareNot
func NotNil(a interface{}) (isNil bool) {
	return web.NotNil(a)
}

// GetConfig get the Appconfig
func GetConfig(returnType, key string, defaultVal interface{}) (interface{}, error) {
	return web.GetConfig(returnType, key, defaultVal)
}

// Str2html Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return web.Str2html(raw)
}

// Htmlquote returns quoted html string.
func Htmlquote(text string) string {
	return web.Htmlquote(text)
}

// Htmlunquote returns unquoted html string.
func Htmlunquote(text string) string {
	return web.Htmlunquote(text)
}

// URLFor returns url string with another registered controller handler with params.
//	usage:
//
//	URLFor(".index")
//	print URLFor("index")
//  router /login
//	print URLFor("login")
//	print URLFor("login", "next","/"")
//  router /profile/:username
//	print UrlFor("profile", ":username","John Doe")
//	result:
//	/
//	/login
//	/login?next=/
//	/user/John%20Doe
//
//  more detail http://beego.vip/docs/mvc/controller/urlbuilding.md
func URLFor(endpoint string, values ...interface{}) string {
	return web.URLFor(endpoint, values...)
}

// AssetsJs returns script tag with src string.
func AssetsJs(text string) template.HTML {
	return web.AssetsJs(text)
}

// AssetsCSS returns stylesheet link tag with src string.
func AssetsCSS(text string) template.HTML {
	text = "<link href=\"" + text + "\" rel=\"stylesheet\" />"

	return template.HTML(text)
}

// ParseForm will parse form values to struct via tag.
func ParseForm(form url.Values, obj interface{}) error {
	return web.ParseForm(form, obj)
}

// RenderForm will render object to form html.
// obj must be a struct pointer.
func RenderForm(obj interface{}) template.HTML {
	return web.RenderForm(obj)
}

// MapGet getting value from map by keys
// usage:
// Data["m"] = M{
//     "a": 1,
//     "1": map[string]float64{
//         "c": 4,
//     },
// }
//
// {{ map_get m "a" }} // return 1
// {{ map_get m 1 "c" }} // return 4
func MapGet(arg1 interface{}, arg2 ...interface{}) (interface{}, error) {
	return web.MapGet(arg1, arg2...)
}
