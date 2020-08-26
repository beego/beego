module github.com/astaxie/beego

require (
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/beego/goyaml2 v0.0.0-20130207012346-5545475820dd
	github.com/beego/x2j v0.0.0-20131220205130-a0352aadc542
	github.com/bradfitz/gomemcache v0.0.0-20180710155616-bc664df96737
	github.com/casbin/casbin v1.7.0
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58
	github.com/couchbase/go-couchbase v0.0.0-20200519150804-63f3cdb75e0d
	github.com/couchbase/gomemcached v0.0.0-20200526233749-ec430f949808 // indirect
	github.com/couchbase/goutils v0.0.0-20180530154633-e865a1461c8a // indirect
	github.com/elastic/go-elasticsearch/v6 v6.8.5
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.1.1
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/hashicorp/golang-lru v0.5.4
	github.com/ledisdb/ledisdb v0.0.0-20200510135210-d35789ec47e6
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.0
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644
	github.com/ssdb/gossdb v0.0.0-20180723034631-88f6b59b84ec
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/goleveldb v0.0.0-20181127023241-353a9fca669c // indirect
	github.com/wendal/errors v0.0.0-20130201093226-f66c77a7882b // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

replace golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 => github.com/golang/crypto v0.0.0-20181127143415-eb0de9b17e85

replace gopkg.in/yaml.v2 v2.2.1 => github.com/go-yaml/yaml v0.0.0-20180328195020-5420a8b6744d

go 1.13
