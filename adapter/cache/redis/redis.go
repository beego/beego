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

// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/beego/beego/v2/client/cache/redis"
//   "github.com/beego/beego/v2/client/cache"
// )
//
//  bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
//
//  more docs http://beego.vip/docs/module/cache.md
package redis

import (
	"github.com/beego/beego/v2/adapter/cache"
	redis2 "github.com/beego/beego/v2/client/cache/redis"
)

// DefaultKey the collection name of redis for cache adapter.
var DefaultKey = "beecacheRedis"

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return cache.CreateNewToOldCacheAdapter(redis2.NewRedisCache())
}

func init() {
	cache.Register("redis", NewRedisCache)
}
