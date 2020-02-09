module github.com/astaxie/beego

require (
	github.com/OwnLocal/goes v1.0.0
	github.com/beego/goyaml2 v0.0.0-20130207012346-5545475820dd
	github.com/beego/x2j v0.0.0-20131220205130-a0352aadc542
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/casbin/casbin v1.9.1
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58
	github.com/couchbase/go-couchbase v0.0.0-20191118180010-b74e3a26e1d7
	github.com/couchbase/gomemcached v0.0.0-20191004160342-7b5da2ec40b2 // indirect
	github.com/couchbase/goutils v0.0.0-20191018232750-b49639060d85 // indirect
	github.com/cupcake/rdb v0.0.0-20161107195141-43ba34106c76 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v2.0.1+incompatible
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/ledisdb v0.0.0-20190202134119-8ceb77e66a92
	github.com/siddontang/rdb v0.0.0-20150307021120-fc89ed2e418d // indirect
	github.com/ssdb/gossdb v0.0.0-20180723034631-88f6b59b84ec
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0
	github.com/wendal/errors v0.0.0-20181209125328-7f31f4b264ec // indirect
	golang.org/x/crypto v0.0.0-20200208060501-ecb85df21340
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

replace golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 => github.com/golang/crypto v0.0.0-20181127143415-eb0de9b17e85

replace gopkg.in/yaml.v2 v2.2.1 => github.com/go-yaml/yaml v0.0.0-20180328195020-5420a8b6744d

go 1.13
