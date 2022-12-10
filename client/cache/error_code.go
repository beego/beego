// Copyright 2021 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"github.com/beego/beego/v2/core/berror"
)

var NilCacheAdapter = berror.DefineCode(4002001, moduleName, "NilCacheAdapter", `
It means that you register cache adapter by pass nil.
A cache adapter is an instance of Cache interface. 
`)

var DuplicateAdapter = berror.DefineCode(4002002, moduleName, "DuplicateAdapter", `
You register two adapter with same name. In beego cache module, one name one adapter.
Once you got this error, please check the error stack, search adapter 
`)

var UnknownAdapter = berror.DefineCode(4002003, moduleName, "UnknownAdapter", `
Unknown adapter, do you forget to register the adapter?
You must register adapter before use it. For example, if you want to use redis implementation, 
you must import the cache/redis package.
`)

var IncrementOverflow = berror.DefineCode(4002004, moduleName, "IncrementOverflow", `
The increment operation will overflow.
`)

var DecrementOverflow = berror.DefineCode(4002005, moduleName, "DecrementOverflow", `
The decrement operation will overflow.
`)

var NotIntegerType = berror.DefineCode(4002006, moduleName, "NotIntegerType", `
The type of value is not (u)int (u)int32 (u)int64. 
When you want to call Incr or Decr function of Cache API, you must confirm that the value's type is one of (u)int (u)int32 (u)int64.
`)

var InvalidFileCacheDirectoryLevelCfg = berror.DefineCode(4002007, moduleName, "InvalidFileCacheDirectoryLevelCfg", `
You pass invalid DirectoryLevel parameter when you try to StartAndGC file cache instance.
This parameter must be a integer, and please check your input.
`)

var InvalidFileCacheEmbedExpiryCfg = berror.DefineCode(4002008, moduleName, "InvalidFileCacheEmbedExpiryCfg", `
You pass invalid EmbedExpiry parameter when you try to StartAndGC file cache instance.
This parameter must be a integer, and please check your input.
`)

var CreateFileCacheDirFailed = berror.DefineCode(4002009, moduleName, "CreateFileCacheDirFailed", `
Beego failed to create file cache directory. There are two cases:
1. You pass invalid CachePath parameter. Please check your input.
2. Beego doesn't have the permission to create this directory. Please check your file mode.
`)

var InvalidFileCachePath = berror.DefineCode(4002010, moduleName, "InvalidFilePath", `
The file path of FileCache is invalid. Please correct the config.
`)

var ReadFileCacheContentFailed = berror.DefineCode(4002011, moduleName, "ReadFileCacheContentFailed", `
Usually you won't got this error. It means that Beego cannot read the data from the file.
You need to check whether the file exist. Sometimes it may be deleted by other processes.
If the file exists, please check the permission that Beego is able to read data from the file.
`)

var InvalidGobEncodedData = berror.DefineCode(4002012, moduleName, "InvalidEncodedData", `
The data is invalid. When you try to decode the invalid data, you got this error.
Please confirm that the data is encoded by GOB correctly.
`)

var GobEncodeDataFailed = berror.DefineCode(4002013, moduleName, "GobEncodeDataFailed", `
Beego could not encode the data to GOB byte array. In general, the data type is invalid. 
For example, GOB doesn't support function type.
Basic types, string, structure, structure pointer are supported. 
`)

var KeyExpired = berror.DefineCode(4002014, moduleName, "KeyExpired", `
Cache key is expired.
You should notice that, a key is expired and then it may be deleted by GC goroutine. 
So when you query a key which may be expired, you may got this code, or KeyNotExist.
`)

var KeyNotExist = berror.DefineCode(4002015, moduleName, "KeyNotExist", `
Key not found.
`)

var MultiGetFailed = berror.DefineCode(4002016, moduleName, "MultiGetFailed", `
Get multiple keys failed. Please check the detail msg to find out the root cause.
`)

var InvalidMemoryCacheCfg = berror.DefineCode(4002017, moduleName, "InvalidMemoryCacheCfg", `
The config is invalid. Please check your input. It must be a json string.
`)

var InvalidMemCacheCfg = berror.DefineCode(4002018, moduleName, "InvalidMemCacheCfg", `
The config is invalid. Please check your input, it must be json string and contains "conn" field.
`)

var InvalidMemCacheValue = berror.DefineCode(4002019, moduleName, "InvalidMemCacheValue", `
The value must be string or byte[], please check your input.
`)

var InvalidRedisCacheCfg = berror.DefineCode(4002020, moduleName, "InvalidRedisCacheCfg", `
The config must be json string, and has "conn" field.
`)

var InvalidSsdbCacheCfg = berror.DefineCode(4002021, moduleName, "InvalidSsdbCacheCfg", `
The config must be json string, and has "conn" field. The value of "conn" field should be "host:port".
"port" must be a valid integer.
`)

var InvalidSsdbCacheValue = berror.DefineCode(4002022, moduleName, "InvalidSsdbCacheValue", `
SSDB cache only accept string value. Please check your input.
`)

var InvalidStoreFunc = berror.DefineCode(4002024, moduleName, "InvalidStoreFunc", `
Invalid storeFunc.
`)

var PersistCacheFailed = berror.DefineCode(4002025, moduleName, "PersistCacheFailed", `
StoreFunc execute failed.
`)

var DeleteFileCacheItemFailed = berror.DefineCode(5002001, moduleName, "DeleteFileCacheItemFailed", `
Beego try to delete file cache item failed. 
Please check whether Beego generated file correctly. 
And then confirm whether this file is already deleted by other processes or other people.
`)

var MemCacheCurdFailed = berror.DefineCode(5002002, moduleName, "MemCacheError", `
When you want to get, put, delete key-value from remote memcache servers, you may get error:
1. You pass invalid servers address, so Beego could not connect to remote server;
2. The servers address is correct, but there is some net issue. Typically there is some firewalls between application and memcache server;
3. Key is invalid. The key's length should be less than 250 and must not contains special characters;
4. The response from memcache server is invalid;
`)

var RedisCacheCurdFailed = berror.DefineCode(5002003, moduleName, "RedisCacheCurdFailed", `
When Beego uses client to send request to redis server, it failed.
1. The server addresses is invalid;
2. Network issue, firewall issue or network is unstable;
3. Client failed to manage connection. In extreme cases, Beego's redis client didn't maintain connections correctly, for example, Beego try to send request via closed connection;
4. The request are huge and redis server spent too much time to process it, and client is timeout;

In general, if you always got this error whatever you do, in most cases, it was caused by network issue. 
You could check your network state, and confirm that firewall rules are correct.
`)

var InvalidConnection = berror.DefineCode(5002004, moduleName, "InvalidConnection", `
The connection is invalid. Please check your connection info, network, firewall.
You could simply uses ping, telnet or write some simple tests to test network.
`)

var DialFailed = berror.DefineCode(5002005, moduleName, "DialFailed", `
When Beego try to dial to remote servers, it failed. Please check your connection info and network state, server state.
`)

var SsdbCacheCurdFailed = berror.DefineCode(5002006, moduleName, "SsdbCacheCurdFailed", `
When you try to use SSDB cache, it failed. There are many cases:
1. servers unavailable;
2. network issue, including network unstable, firewall;
3. connection issue;
4. request are huge and servers spent too much time to process it, got timeout;
`)

var SsdbBadResponse = berror.DefineCode(5002007, moduleName, "SsdbBadResponse", `
The reponse from SSDB server is invalid. 
Usually it indicates something wrong on server side.
`)

var (
	ErrKeyExpired  = berror.Error(KeyExpired, "the key is expired")
	ErrKeyNotExist = berror.Error(KeyNotExist, "the key isn't exist")
)
