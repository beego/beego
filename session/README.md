session
==============

session is a Go session manager. It can use many session providers. Just like the `database/sql` and `database/sql/driver`.

## How to install?

	go get github.com/astaxie/beego/session


## What providers are supported?

As of now this session manager support memory, file, Redis and MySQL.


## How to use it?

First you must import it

	import (
		"github.com/astaxie/beego/session"
	)

Then in you web app init the global session manager
	
	var globalSessions *session.Manager

* Use **memory** as provider:

		func init() {
			globalSessions, _ = session.NewManager("memory", "gosessionid", 3600,"")
			go globalSessions.GC()
		}

* Use **file** as provider, the last param is the path where you want file to be stored:

		func init() {
			globalSessions, _ = session.NewManager("file", "gosessionid", 3600, "./tmp")
			go globalSessions.GC()
		}

* Use **Redis** as provider, the last param is the Redis conn address:

		func init() {
			globalSessions, _ = session.NewManager("redis", "gosessionid", 3600, "127.0.0.1:6379")
			go globalSessions.GC()
		}
		
* Use **MySQL** as provider, the last param is the DSN, learn more from [mysql](https://github.com/Go-SQL-Driver/MySQL#dsn-data-source-name): 

		func init() {
			globalSessions, _ = session.NewManager(
				"mysql", "gosessionid", 3600, "username:password@protocol(address)/dbname?param=value")
			go globalSessions.GC()
		}

Finally in the handlerfunc you can use it like this

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


## How to write own provider?

When you develop a web app, maybe you want to write own provider because you must meet the requirements.

Writing a provider is easy. You only need to define two struct types 
(Session and Provider), which satisfy the interface definition. 
Maybe you will find the **memory** provider as good example.

	type SessionStore interface {
		Set(key, value interface{}) error // set session value
		Get(key interface{}) interface{}  // get session value
		Delete(key interface{}) error     // delete session value
		SessionID() string                // return current sessionID
		SessionRelease()                  // release the resource
	}
	
	type Provider interface {
		SessionInit(maxlifetime int64, savePath string) error
		SessionRead(sid string) (SessionStore, error)
		SessionDestroy(sid string) error
		SessionGC()
	}


## LICENSE

BSD License http://creativecommons.org/licenses/BSD/
