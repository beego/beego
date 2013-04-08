sessionmanager
==============

sessionmanager is a golang session manager. It can use session for many providers.Just like the `database/sql` and `database/sql/driver`.

##How to install

	go get github.com/astaxie/beego/session


##how many providers support
Now this sessionmanager support memory/file/redis/mysql



##How do we use it?

first you must import it


	import (
		"github.com/astaxie/beego/session"
	)

then in you web app init the globalsession manager
	
	var globalSessions *session.Manager

use memory as providers:

	func init() {
		globalSessions, _ = session.NewManager("memory", "gosessionid", 3600,"")
		go globalSessions.GC()
	}

use mysql as providers,the last param is the DNS, learn more from [mysql](https://github.com/Go-SQL-Driver/MySQL#dsn-data-source-name): 

	func init() {
		globalSessions, _ = session.NewManager("mysql", "gosessionid", 3600,"username:password@protocol(address)/dbname?param=value")
		go globalSessions.GC()
	}

use file as providers,the last param is the path where to store the file:

	func init() {
		globalSessions, _ = session.NewManager("file", "gosessionid", 3600,"./tmp")
		go globalSessions.GC()
	}

use redis as providers,the last param is the redis's conn address:

	func init() {
		globalSessions, _ = session.NewManager("redis", "gosessionid", 3600,"127.0.0.1:6379")
		go globalSessions.GC()
	}

at last in the handlerfunc you can use it like this

	func login(w http.ResponseWriter, r *http.Request) {
		sess := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease()
		username := sess.Get("username")
		fmt.Println(username)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("login.gtpl")
			t.Execute(w, nil)
		} else {
			fmt.Println("username:", r.Form["username"])
			sess.Set("username", r.Form["username"])
			fmt.Println("password:", r.Form["password"])
		}
	}
	


##How to write own provider
When we develop a web app, maybe you want to write a provider because you must meet the requirements.

Write a provider is so easy. You only define two struct type(Session and Provider),which satisfy the interface definition.Maybe The memory provider is a good example for you.

	type SessionStore interface {
		Set(key, value interface{}) error //set session value
		Get(key interface{}) interface{}  //get session value
		Delete(key interface{}) error     //delete session value
		SessionID() string                //back current sessionID
		SessionRelease()                  // release the resource
	}
	
	type Provider interface {
		SessionInit(maxlifetime int64, savePath string) error
		SessionRead(sid string) (SessionStore, error)
		SessionDestroy(sid string) error
		SessionGC()
	}

##LICENSE

BSD License http://creativecommons.org/licenses/BSD/