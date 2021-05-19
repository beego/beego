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

package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beego/beego/v2/client/httplib"
)

func TestSimpleCondition_MatchPath(t *testing.T) {
	sc := NewSimpleCondition("/abc/s")
	res := sc.Match(context.Background(), httplib.Get("http://localhost:8080/abc/s"))
	assert.True(t, res)
}

func TestSimpleCondition_MatchQuery(t *testing.T) {
	k, v := "my-key", "my-value"
	sc := NewSimpleCondition("/abc/s")
	res := sc.Match(context.Background(), httplib.Get("http://localhost:8080/abc/s?my-key=my-value"))
	assert.True(t, res)

	sc = NewSimpleCondition("/abc/s", WithQuery(k, v))
	res = sc.Match(context.Background(), httplib.Get("http://localhost:8080/abc/s?my-key=my-value"))
	assert.True(t, res)

	res = sc.Match(context.Background(), httplib.Get("http://localhost:8080/abc/s?my-key=my-valuesss"))
	assert.False(t, res)

	res = sc.Match(context.Background(), httplib.Get("http://localhost:8080/abc/s?my-key-a=my-value"))
	assert.False(t, res)

	res = sc.Match(context.Background(), httplib.Get("http://localhost:8080/abc/s?my-key=my-value&abc=hello"))
	assert.True(t, res)
}

func TestSimpleCondition_MatchHeader(t *testing.T) {
	k, v := "my-header", "my-header-value"
	sc := NewSimpleCondition("/abc/s")
	req := httplib.Get("http://localhost:8080/abc/s")
	assert.True(t, sc.Match(context.Background(), req))

	req = httplib.Get("http://localhost:8080/abc/s")
	req.Header(k, v)
	assert.True(t, sc.Match(context.Background(), req))

	sc = NewSimpleCondition("/abc/s", WithHeader(k, v))
	req.Header(k, v)
	assert.True(t, sc.Match(context.Background(), req))

	req.Header(k, "invalid")
	assert.False(t, sc.Match(context.Background(), req))
}

func TestSimpleCondition_MatchBodyField(t *testing.T) {
	sc := NewSimpleCondition("/abc/s")
	req := httplib.Post("http://localhost:8080/abc/s")

	assert.True(t, sc.Match(context.Background(), req))

	req.Body(`{
    "body-field": 123
}`)
	assert.True(t, sc.Match(context.Background(), req))

	k := "body-field"
	v := float64(123)
	sc = NewSimpleCondition("/abc/s", WithJsonBodyFields(k, v))
	assert.True(t, sc.Match(context.Background(), req))

	sc = NewSimpleCondition("/abc/s", WithJsonBodyFields(k, v))
	req.Body(`{
    "body-field": abc
}`)
	assert.False(t, sc.Match(context.Background(), req))

	sc = NewSimpleCondition("/abc/s", WithJsonBodyFields("body-field", "abc"))
	req.Body(`{
    "body-field": "abc"
}`)
	assert.True(t, sc.Match(context.Background(), req))
}

func TestSimpleCondition_Match(t *testing.T) {
	sc := NewSimpleCondition("/abc/s")
	req := httplib.Post("http://localhost:8080/abc/s")

	assert.True(t, sc.Match(context.Background(), req))

	sc = NewSimpleCondition("/abc/s", WithMethod("POST"))
	assert.True(t, sc.Match(context.Background(), req))

	sc = NewSimpleCondition("/abc/s", WithMethod("GET"))
	assert.False(t, sc.Match(context.Background(), req))
}

func TestSimpleCondition_MatchPathReg(t *testing.T) {
	sc := NewSimpleCondition("", WithPathReg(`\/abc\/.*`))
	req := httplib.Post("http://localhost:8080/abc/s")
	assert.True(t, sc.Match(context.Background(), req))

	req = httplib.Post("http://localhost:8080/abcd/s")
	assert.False(t, sc.Match(context.Background(), req))
}
