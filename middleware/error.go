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

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"strconv"
)

var (
	AppName string
	VERSION string
)
var tpl = `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>beego application error</title>
    <style>
        html, body, body * {padding: 0; margin: 0;}
        #header {background:#ffd; border-bottom:solid 2px #A31515; padding: 20px 10px;}
        #header h2{ }
        #footer {border-top:solid 1px #aaa; padding: 5px 10px; font-size: 12px; color:green;}
        #content {padding: 5px;}
        #content .stack b{ font-size: 13px; color: red;}
        #content .stack pre{padding-left: 10px;}
        table {}
        td.t {text-align: right; padding-right: 5px; color: #888;}
    </style>
    <script type="text/javascript">
    </script>
</head>
<body>
    <div id="header">
        <h2>{{.AppError}}</h2>
    </div>
    <div id="content">
        <table>
            <tr>
                <td class="t">Request Method: </td><td>{{.RequestMethod}}</td>
            </tr>
            <tr>
                <td class="t">Request URL: </td><td>{{.RequestURL}}</td>
            </tr>
            <tr>
                <td class="t">RemoteAddr: </td><td>{{.RemoteAddr }}</td>
            </tr>
        </table>
        <div class="stack">
            <b>Stack</b>
            <pre>{{.Stack}}</pre>
        </div>
    </div>
    <div id="footer">
        <p>beego {{ .BeegoVersion }} (beego framework)</p>
        <p>golang version: {{.GoVersion}}</p>
    </div>
</body>
</html>
`

// render default application error page with error and stack string.
func ShowErr(err interface{}, rw http.ResponseWriter, r *http.Request, Stack string) {
	t, _ := template.New("beegoerrortemp").Parse(tpl)
	data := make(map[string]string)
	data["AppError"] = AppName + ":" + fmt.Sprint(err)
	data["RequestMethod"] = r.Method
	data["RequestURL"] = r.RequestURI
	data["RemoteAddr"] = r.RemoteAddr
	data["Stack"] = Stack
	data["BeegoVersion"] = VERSION
	data["GoVersion"] = runtime.Version()
	rw.WriteHeader(500)
	t.Execute(rw, data)
}

var errtpl = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
		<title>{{.Title}}</title>
		<style type="text/css">
			* {
				margin:0;
				padding:0;
			}

			body {
				background-color:#EFEFEF;
				font: .9em "Lucida Sans Unicode", "Lucida Grande", sans-serif;
			}

			#wrapper{
				width:600px;
				margin:40px auto 0;
				text-align:center;
				-moz-box-shadow: 5px 5px 10px rgba(0,0,0,0.3);
				-webkit-box-shadow: 5px 5px 10px rgba(0,0,0,0.3);
				box-shadow: 5px 5px 10px rgba(0,0,0,0.3);
			}

			#wrapper h1{
				color:#FFF;
				text-align:center;
				margin-bottom:20px;
			}

			#wrapper a{
				display:block;
				font-size:.9em;
				padding-top:20px;
				color:#FFF;
				text-decoration:none;
				text-align:center;
			}

			#container {
				width:600px;
				padding-bottom:15px;
				background-color:#FFFFFF;
			}

			.navtop{
				height:40px;
				background-color:#24B2EB;
				padding:13px;
			}

			.content {
				padding:10px 10px 25px;
				background: #FFFFFF;
				margin:;
				color:#333;
			}

			a.button{
				color:white;
				padding:15px 20px;
				text-shadow:1px 1px 0 #00A5FF;
				font-weight:bold;
				text-align:center;
				border:1px solid #24B2EB;
				margin:0px 200px;
				clear:both;
				background-color: #24B2EB;
				border-radius:100px;
				-moz-border-radius:100px;
				-webkit-border-radius:100px;
			}

			a.button:hover{
				text-decoration:none;
				background-color: #24B2EB;
			}

		</style>
	</head>
	<body>
		<div id="wrapper">
			<div id="container">
				<div class="navtop">
					<h1>{{.Title}}</h1>
				</div>
				<div id="content">
					{{.Content}}
					<a href="/" title="Home" class="button">Go Home</a><br />

					<br>Powered by beego {{.BeegoVersion}}
				</div>
			</div>
		</div>
	</body>
</html>
`

// map of http handlers for each error string.
var ErrorMaps map[string]http.HandlerFunc

func init() {
	ErrorMaps = make(map[string]http.HandlerFunc)
}

// show 404 notfound error.
func NotFound(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Page Not Found"
	data["Content"] = template.HTML("<br>The page you have requested has flown the coop." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>The page has moved" +
		"<br>The page no longer exists" +
		"<br>You were looking for your puppy and got lost" +
		"<br>You like 404 pages" +
		"</ul>")
	data["BeegoVersion"] = VERSION
	//rw.WriteHeader(http.StatusNotFound)
	t.Execute(rw, data)
}

// show 401 unauthorized error.
func Unauthorized(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Unauthorized"
	data["Content"] = template.HTML("<br>The page you have requested can't be authorized." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>The credentials you supplied are incorrect" +
		"<br>There are errors in the website address" +
		"</ul>")
	data["BeegoVersion"] = VERSION
	//rw.WriteHeader(http.StatusUnauthorized)
	t.Execute(rw, data)
}

// show 403 forbidden error.
func Forbidden(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Forbidden"
	data["Content"] = template.HTML("<br>The page you have requested is forbidden." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>Your address may be blocked" +
		"<br>The site may be disabled" +
		"<br>You need to log in" +
		"</ul>")
	data["BeegoVersion"] = VERSION
	//rw.WriteHeader(http.StatusForbidden)
	t.Execute(rw, data)
}

// show 503 service unavailable error.
func ServiceUnavailable(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Service Unavailable"
	data["Content"] = template.HTML("<br>The page you have requested is unavailable." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br><br>The page is overloaded" +
		"<br>Please try again later." +
		"</ul>")
	data["BeegoVersion"] = VERSION
	//rw.WriteHeader(http.StatusServiceUnavailable)
	t.Execute(rw, data)
}

// show 500 internal server error.
func InternalServerError(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Internal Server Error"
	data["Content"] = template.HTML("<br>The page you have requested is down right now." +
		"<br><br><ul>" +
		"<br>Please try again later and report the error to the website administrator" +
		"<br></ul>")
	data["BeegoVersion"] = VERSION
	//rw.WriteHeader(http.StatusInternalServerError)
	t.Execute(rw, data)
}

// show 500 internal error with simple text string.
func SimpleServerError(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// add http handler for given error string.
func Errorhandler(err string, h http.HandlerFunc) {
	ErrorMaps[err] = h
}

// register default error http handlers, 404,401,403,500 and 503.
func RegisterErrorHandler() {
	if _, ok := ErrorMaps["404"]; !ok {
		ErrorMaps["404"] = NotFound
	}

	if _, ok := ErrorMaps["401"]; !ok {
		ErrorMaps["401"] = Unauthorized
	}

	if _, ok := ErrorMaps["403"]; !ok {
		ErrorMaps["403"] = Forbidden
	}

	if _, ok := ErrorMaps["503"]; !ok {
		ErrorMaps["503"] = ServiceUnavailable
	}

	if _, ok := ErrorMaps["500"]; !ok {
		ErrorMaps["500"] = InternalServerError
	}
}

// show error string as simple text message.
// if error string is empty, show 500 error as default.
func Exception(errcode string, w http.ResponseWriter, r *http.Request, msg string) {
	if h, ok := ErrorMaps[errcode]; ok {
		isint, err := strconv.Atoi(errcode)
		if err != nil {
			isint = 500
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(isint)
		h(w, r)
		return
	} else {
		isint, err := strconv.Atoi(errcode)
		if err != nil {
			isint = 500
		}
		if isint == 400 {
			msg = "404 page not found"
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(isint)
		fmt.Fprintln(w, msg)
		return
	}
}
