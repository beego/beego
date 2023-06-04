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
	"net/http"

	"github.com/beego/beego/v2/adapter/context"
	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
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

// ErrorMaps holds map of http handlers for each error string.
// there is 10 kinds default error(40x and 50x)
var ErrorMaps = web.ErrorMaps

// ErrorHandler registers http.HandlerFunc to each http err code string.
// usage:
//
//	beego.ErrorHandler("404",NotFound)
//	beego.ErrorHandler("500",InternalServerError)
func ErrorHandler(code string, h http.HandlerFunc) *App {
	return (*App)(web.ErrorHandler(code, h))
}

// ErrorController registers ControllerInterface to each http err code string.
// usage:
//
//	beego.ErrorController(&controllers.ErrorController{})
func ErrorController(c ControllerInterface) *App {
	return (*App)(web.ErrorController(c))
}

// Exception Write HttpStatus with errCode and Exec error handler if exist.
func Exception(errCode uint64, ctx *context.Context) {
	web.Exception(errCode, (*beecontext.Context)(ctx))
}
