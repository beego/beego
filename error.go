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
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/utils"
)

const (
	errorTypeHandler = iota
	errorTypeController
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
func showErr(err interface{}, ctx *context.Context, Stack string) {
	t, _ := template.New("beegoerrortemp").Parse(tpl)
	data := make(map[string]string)
	data["AppError"] = AppName + ":" + fmt.Sprint(err)
	data["RequestMethod"] = ctx.Input.Method()
	data["RequestURL"] = ctx.Input.Uri()
	data["RemoteAddr"] = ctx.Input.IP()
	data["Stack"] = Stack
	data["BeegoVersion"] = VERSION
	data["GoVersion"] = runtime.Version()
	ctx.ResponseWriter.WriteHeader(500)
	t.Execute(ctx.ResponseWriter, data)
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

type errorInfo struct {
	controllerType reflect.Type
	handler        http.HandlerFunc
	method         string
	errorType      int
}

// map of http handlers for each error string.
var ErrorMaps map[string]*errorInfo

func init() {
	ErrorMaps = make(map[string]*errorInfo)
}

// show 401 unauthorized error.
func unauthorized(rw http.ResponseWriter, r *http.Request) {
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
	t.Execute(rw, data)
}

// show 402 Payment Required
func paymentRequired(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Payment Required"
	data["Content"] = template.HTML("<br>The page you have requested Payment Required." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>The credentials you supplied are incorrect" +
		"<br>There are errors in the website address" +
		"</ul>")
	data["BeegoVersion"] = VERSION
	t.Execute(rw, data)
}

// show 403 forbidden error.
func forbidden(rw http.ResponseWriter, r *http.Request) {
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
	t.Execute(rw, data)
}

// show 404 notfound error.
func notFound(rw http.ResponseWriter, r *http.Request) {
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
	t.Execute(rw, data)
}

// show 405 Method Not Allowed
func methodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Method Not Allowed"
	data["Content"] = template.HTML("<br>The method you have requested Not Allowed." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>The method specified in the Request-Line is not allowed for the resource identified by the Request-URI" +
		"<br>The response MUST include an Allow header containing a list of valid methods for the requested resource." +
		"</ul>")
	data["BeegoVersion"] = VERSION
	t.Execute(rw, data)
}

// show 500 internal server error.
func internalServerError(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Internal Server Error"
	data["Content"] = template.HTML("<br>The page you have requested is down right now." +
		"<br><br><ul>" +
		"<br>Please try again later and report the error to the website administrator" +
		"<br></ul>")
	data["BeegoVersion"] = VERSION
	t.Execute(rw, data)
}

// show 501 Not Implemented.
func notImplemented(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Not Implemented"
	data["Content"] = template.HTML("<br>The page you have requested is Not Implemented." +
		"<br><br><ul>" +
		"<br>Please try again later and report the error to the website administrator" +
		"<br></ul>")
	data["BeegoVersion"] = VERSION
	t.Execute(rw, data)
}

// show 502 Bad Gateway.
func badGateway(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Bad Gateway"
	data["Content"] = template.HTML("<br>The page you have requested is down right now." +
		"<br><br><ul>" +
		"<br>The server, while acting as a gateway or proxy, received an invalid response from the upstream server it accessed in attempting to fulfill the request." +
		"<br>Please try again later and report the error to the website administrator" +
		"<br></ul>")
	data["BeegoVersion"] = VERSION
	t.Execute(rw, data)
}

// show 503 service unavailable error.
func serviceUnavailable(rw http.ResponseWriter, r *http.Request) {
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
	t.Execute(rw, data)
}

// show 504 Gateway Timeout.
func gatewayTimeout(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("beegoerrortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Gateway Timeout"
	data["Content"] = template.HTML("<br>The page you have requested is unavailable." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br><br>The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server specified by the URI." +
		"<br>Please try again later." +
		"</ul>")
	data["BeegoVersion"] = VERSION
	t.Execute(rw, data)
}

// register default error http handlers, 404,401,403,500 and 503.
func registerDefaultErrorHandler() {
	if _, ok := ErrorMaps["401"]; !ok {
		Errorhandler("401", unauthorized)
	}

	if _, ok := ErrorMaps["402"]; !ok {
		Errorhandler("402", paymentRequired)
	}

	if _, ok := ErrorMaps["403"]; !ok {
		Errorhandler("403", forbidden)
	}

	if _, ok := ErrorMaps["404"]; !ok {
		Errorhandler("404", notFound)
	}

	if _, ok := ErrorMaps["405"]; !ok {
		Errorhandler("405", methodNotAllowed)
	}

	if _, ok := ErrorMaps["500"]; !ok {
		Errorhandler("500", internalServerError)
	}
	if _, ok := ErrorMaps["501"]; !ok {
		Errorhandler("501", notImplemented)
	}
	if _, ok := ErrorMaps["502"]; !ok {
		Errorhandler("502", badGateway)
	}

	if _, ok := ErrorMaps["503"]; !ok {
		Errorhandler("503", serviceUnavailable)
	}

	if _, ok := ErrorMaps["504"]; !ok {
		Errorhandler("504", gatewayTimeout)
	}
}

// ErrorHandler registers http.HandlerFunc to each http err code string.
// usage:
// 	beego.ErrorHandler("404",NotFound)
//	beego.ErrorHandler("500",InternalServerError)
func Errorhandler(code string, h http.HandlerFunc) *App {
	errinfo := &errorInfo{}
	errinfo.errorType = errorTypeHandler
	errinfo.handler = h
	errinfo.method = code
	ErrorMaps[code] = errinfo
	return BeeApp
}

// ErrorController registers ControllerInterface to each http err code string.
// usage:
// 	beego.ErrorHandler(&controllers.ErrorController{})
func ErrorController(c ControllerInterface) *App {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	for i := 0; i < rt.NumMethod(); i++ {
		if !utils.InSlice(rt.Method(i).Name, exceptMethod) && strings.HasPrefix(rt.Method(i).Name, "Error") {
			errinfo := &errorInfo{}
			errinfo.errorType = errorTypeController
			errinfo.controllerType = ct
			errinfo.method = rt.Method(i).Name
			errname := strings.TrimPrefix(rt.Method(i).Name, "Error")
			ErrorMaps[errname] = errinfo
		}
	}
	return BeeApp
}

// show error string as simple text message.
// if error string is empty, show 500 error as default.
func exception(errcode string, ctx *context.Context) {
	code, err := strconv.Atoi(errcode)
	if err != nil {
		code = 503
	}
	if h, ok := ErrorMaps[errcode]; ok {
		executeError(h, ctx, code)
		return
	} else if h, ok := ErrorMaps["503"]; ok {
		executeError(h, ctx, code)
		return
	} else {
		ctx.ResponseWriter.WriteHeader(code)
		ctx.WriteString(errcode)
	}
}

func executeError(err *errorInfo, ctx *context.Context, code int) {
	if err.errorType == errorTypeHandler {
		ctx.ResponseWriter.WriteHeader(code)
		err.handler(ctx.ResponseWriter, ctx.Request)
		return
	}
	if err.errorType == errorTypeController {
		ctx.Output.SetStatus(code)
		//Invoke the request handler
		vc := reflect.New(err.controllerType)
		execController, ok := vc.Interface().(ControllerInterface)
		if !ok {
			panic("controller is not ControllerInterface")
		}
		//call the controller init function
		execController.Init(ctx, err.controllerType.Name(), err.method, vc.Interface())

		//call prepare function
		execController.Prepare()

		execController.URLMapping()

		in := make([]reflect.Value, 0)
		method := vc.MethodByName(err.method)
		method.Call(in)

		//render template
		if AutoRender {
			if err := execController.Render(); err != nil {
				panic(err)
			}
		}

		// finish all runrouter. release resource
		execController.Finish()
	}
}
