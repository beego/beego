// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package memcache for cache provider
//
// depend on github.com/bradfitz/gomemcache/memcache
//
// go install github.com/bradfitz/gomemcache/memcache
//
// Usage:
// import(
//   _ "github.com/beego/beego/v2/cache/memcache"
//   "github.com/beego/beego/v2/cache"
// )
//
//  bm, err := cache.NewCache("memcache", `{"conn":"127.0.0.1:11211"}`)
//
//  more docs http://beego.me/docs/module/cache.md
package memcache

import (
	"github.com/beego/beego/v2/adapter/cache"
	"github.com/beego/beego/v2/client/cache/memcache"
)

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return cache.CreateNewToOldCacheAdapter(memcache.NewMemCache())
}

func init() {
	cache.Register("memcache", NewMemCache)
}
