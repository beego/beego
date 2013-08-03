## 方便的http客户端
我们经常会使用Go来请求其他API应用，例如你使用beego开发了一个RESTFul的API应用，那么如果来请求呢？当然可以使用`http.Client`来实现，但是需要自己来操作很多步骤，自己需要考虑很多东西，所以我就基于net下的一些包实现了这个简便的http客户端工具。

该工具的主要特点：

- 链式操作
- 超时控制
- 方便的解析
- 可控的debug

## 例子
我们上次开发的RESTful应用，最后我写过如何通过curl来进行测试，那么下面一一对每个操作如何用httplib来操作进行展示

- 添加一个对象：

	`curl -X POST -d '{"Score":1337,"PlayerName":"Sean Plott"}' http://127.0.0.1:8080/object`
	
	返回一个相应的objectID:astaxie1373349756660423900
	
		str,err:=beego.Post("http://127.0.0.1:8080/object").Body(`{"Score":1337,"PlayerName":"Sean Plott"}`).String()
		if err != nil{
			println(err)
		}
	
- 查询一个对象

	`curl -X GET http://127.0.0.1:8080/object/astaxie1373349756660423900`
	
		var object Obeject
		err:=beego.Get("http://127.0.0.1:8080/object/astaxie1373349756660423900").ToJson(&object)
		if err != nil{
			println(err)
		}
	
- 查询全部的对象	

	`curl -X GET http://127.0.0.1:8080/object`
	
		var objects []Object
		err:=beego.Get("http://127.0.0.1:8080/object").ToJson(&objects)
		if err != nil{
			println(err)
		}

- 更新一个对象

	`curl -X PUT -d '{"Score":10000}'http://127.0.0.1:8080/object/astaxie1373349756660423900`
	
		str,err:=beego.Put("http://127.0.0.1:8080/object/astaxie1373349756660423900").Body(`{"Score":10000}`).String()
		if err != nil{
			println(err)
		}

- 删除一个对象

	`curl -X DELETE http://127.0.0.1:8080/object/astaxie1373349756660423900`
	
		str,er:=beego.Delete("http://127.0.0.1:8080/object/astaxie1373349756660423900").String()
		if err != nil{
			println(err)
		}

## 开启调试模式
用户可以开启调试打印request信息，默认是关闭模式

	beego.Post(url).Debug(true)

## ToFile、ToXML、ToJson
上面我演示了Json的解析，其实还有直接保存为文件的ToFile操作，解析XML的ToXML操作


## 设置链接超时和读写超时
默认都设置为60秒，用户可以通过函数来设置相应的超时时间

	beego.Get(url).SetTimeout(100*time.Second,100*time.Second)


更加详细的请参考[API接口](http://gowalker.org/github.com/astaxie/beego)