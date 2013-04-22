## cache
cache is a golang cache manager. It can use cache for many adapters. The repo is inspired by `database/sql` .

##How to install

	go get github.com/astaxie/beego/cache

	
##how many adapter support

Now this cache support memory/redis/memcache	
	
## how to use it
first you must import it


	import (
		"github.com/astaxie/beego/cache"
	)

then init an Cache(memory adapter)

	bm, err := NewCache("memory", `{"interval":60}`)	

use it like this:	
	
	bm.Put("astaxie", 1, 10)
	bm.Get("astaxie")
	bm.IsExist("astaxie")
	bm.Delete("astaxie")
	
## memory adapter
memory adapter config like this:

	{"interval":60}

interval means the gc time. The cache will every interval time to check wheather have item expired.	

## memcache adapter
memory adapter use the vitess's [memcache](code.google.com/p/vitess/go/memcache) client.

the config like this:

	{"conn":"127.0.0.1:11211"}


## redis	 adapter
redis adapter use the [redigo](github.com/garyburd/redigo/redis) client.

the config like this:

	{"conn":":6039"}