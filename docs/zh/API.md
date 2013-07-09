# API应用开发入门
Go是非常适合用来开发API应用的，而且我认为也是Go相对于其他动态语言的最大优势应用。beego在开发API应用方面提供了非常强大和快速的工具，方便用户快速的建立API应用原型，专心业务逻辑就行了。


## 快速建立原型
bee快速开发工具提供了一个API应用建立的工具，在gopath/src下的任意目录执行如下命令就可以快速的建立一个API应用：

`bee api beeapi`

## 应用的目录结构
应用的目录结构如下所示：

 ```
├── conf
│   └── app.conf
├── controllers
│   └── default.go
├── models
│    └── object.go
└── main.go	
```

## 源码解析

- app.conf里面主要针对API的配置如下：

autorender = false  //API应用不需要模板渲染，所以关闭自动渲染

copyrequestbody = true  //RESTFul应用发送信息的时候是raw body，而不是普通的form表单，所以需要额外的读取body信息

- main.go文件主要针对RESTFul的路由注册

`beego.RESTRouter("/object", &controllers.ObejctController{})`

这个路由可以匹配如下的规则

<table>
<tr>
	<th>URL</th>					<th>HTTP Verb</th>				<th>Functionality</th>
</tr>
<tr>	
	<td>/object</td>				<td>POST</td>					<td>Creating Objects</td>
</tr>
<tr>	
	<td>/object/objectId</td>	<td>GET</td>					<td>Retrieving Objects</td>
</tr>
<tr>	
	<td>/object/objectId</td>	<td>PUT</td>						<td>Updating Objects</td>
</tr>
<tr>	
	<td>/object</td>				<td>GET</td>						<td>Queries</td>
</tr>
<tr>	
	<td>/object/objectId</td>	<td>DELETE</td>					<td>Deleting Objects</td>
</tr>
</table>

- ObejctController实现了对应的方法：

```
type ObejctController struct {
	beego.Controller
}

func (this *ObejctController) Post(){
	
}

func (this *ObejctController) Get(){
	
}

func (this *ObejctController) Put(){
	
}

func (this *ObejctController) Delete(){
	
}
```

- models里面实现了对应操作对象的增删改取等操作

## 测试

- 添加一个对象：

	`curl -X POST -d '{"Score":1337,"PlayerName":"Sean Plott"}' http://127.0.0.1:8080/object`
	
	返回一个相应的objectID:astaxie1373349756660423900
	
- 查询一个对象

	`curl -X GET http://127.0.0.1:8080/object/astaxie1373349756660423900`
	
- 查询全部的对象	

	`curl -X GET http://127.0.0.1:8080/object`

- 更新一个对象

	`curl -X PUT -d '{"Score":10000}'http://127.0.0.1:8080/object/astaxie1373349756660423900`

- 删除一个对象

	`curl -X DELETE http://127.0.0.1:8080/object/astaxie1373349756660423900`
