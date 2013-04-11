package beego

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime"
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

func ShowErr(err interface{}, rw http.ResponseWriter, r *http.Request, Stack string) {
	t, err := template.New("beegoerrortemp").Parse(tpl)
	data := make(map[string]string)
	data["AppError"] = AppName + ":" + fmt.Sprint(err)
	data["RequestMethod"] = r.Method
	data["RequestURL"] = r.RequestURI
	data["RemoteAddr"] = r.RemoteAddr
	data["Stack"] = Stack
	data["BeegoVersion"] = VERSION
	data["GoVersion"] = runtime.Version()
	t.Execute(rw, data)
}
