// Copyright 2020 beego
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

package httplib

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

// BeegoHTTPSettings is the http.Client setting
type BeegoHTTPSettings struct {
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	EnableCookie     bool
	Gzip             bool
	Retries          int // if set to -1 means will retry forever
	RetryDelay       time.Duration
	FilterChains     []FilterChain
}

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// SetDefaultSetting overwrites default settings
// Keep in mind that when you invoke the SetDefaultSetting
// some methods invoked before SetDefaultSetting
func SetDefaultSetting(setting BeegoHTTPSettings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultSetting = setting
}

var defaultSetting = BeegoHTTPSettings{
	UserAgent:        "beegoServer",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	FilterChains:     make([]FilterChain, 0, 4),
}

var (
	defaultCookieJar http.CookieJar
	settingMutex     sync.Mutex
)

// AddDefaultFilter add a new filter into defaultSetting
// Be careful about using this method if you invoke SetDefaultSetting somewhere
func AddDefaultFilter(fc FilterChain) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	if defaultSetting.FilterChains == nil {
		defaultSetting.FilterChains = make([]FilterChain, 0, 4)
	}
	defaultSetting.FilterChains = append(defaultSetting.FilterChains, fc)
}
