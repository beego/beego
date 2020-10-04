package ssdb

import (
	"github.com/astaxie/beego/pkg/adapter/cache"
	ssdb2 "github.com/astaxie/beego/pkg/client/cache/ssdb"
)

// NewSsdbCache create new ssdb adapter.
func NewSsdbCache() cache.Cache {
	return cache.CreateNewToOldCacheAdapter(ssdb2.NewSsdbCache())
}

func init() {
	cache.Register("ssdb", NewSsdbCache)
}
