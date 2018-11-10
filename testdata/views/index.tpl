<!DOCTYPE html>
<html>
  <head>
    <title>beego welcome template</title>
  </head>
  <body>

	{{template "block"}}
	{{template "header"}}
	{{template "blocks/block.tpl"}}

	<h2>{{ .Title }}</h2>
	<p> This is SomeVar: {{ .SomeVar }}</p>
  </body>
</html>
