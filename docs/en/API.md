# Getting start with API application development
Go is very good for developing API applications which I think is the biggest strength compare to other dynamic languages. Beego provides powerful and quick setup tool for developing API applications, which gives you more focus on business logic.


## Quick setup
bee can setup a API application very quick by executing commands under any `$GOPATH/src`. 

`bee api beeapi`

## Application directory structure

 ```
├── conf
│   └── app.conf
├── controllers
│   └── default.go
├── models
│    └── object.go
└── main.go	
```

## Source code explanation

- app.conf has following configuration options for your API applications:

	- autorender = false  // Disable auto-render since API applications don't need.
	- copyrequestbody = true  // RESTFul applications sends raw body instead of form, so we need to read body specifically.

- main.go is for registering routers of RESTFul.

		beego.RESTRouter("/object", &controllers.ObejctController{})

Match rules as follows:

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

- ObejctController implemented corresponding methods:

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

- models implemented corresponding object operation for adding, deleting, updating and getting.

## Test

- Add a new object:
	
		curl -X POST -d '{"Score":1337,"PlayerName":"Sean Plott"}' http://127.0.0.1:8080/object
		
	Returns a corresponding objectID:astaxie1373349756660423900
	
- Query a object:

	`curl -X GET http://127.0.0.1:8080/object/astaxie1373349756660423900`
	
- Query all objects:

	`curl -X GET http://127.0.0.1:8080/object`

- Update a object:

	`curl -X PUT -d '{"Score":10000}'http://127.0.0.1:8080/object/astaxie1373349756660423900`

- Delete a object:

	`curl -X DELETE http://127.0.0.1:8080/object/astaxie1373349756660423900`
