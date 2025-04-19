// Copyright 2016 beego Author. All Rights Reserved.
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

package context

import (
	"github.com/beego/beego/v2/server/web/session"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test concurrency safety
func TestEventStreamResp_Concurrency(t *testing.T) {
	req, err := http.NewRequest("GET", "/events", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp := httptest.NewRecorder()
	beegoResp := &Response{ResponseWriter: resp}

	ctxInstance := &Context{
		Request:        req,
		ResponseWriter: beegoResp,
	}

	eventCh := ctxInstance.EventStreamResp()

	var wg sync.WaitGroup
	messages := []string{"msg1", "msg2", "msg3", "msg4", "msg5"}

	for _, msg := range messages {
		wg.Add(1)
		go func(m string) {
			defer wg.Done()
			eventCh <- []byte(m)
		}(msg)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	received := resp.Body.String()
	for _, msg := range messages {
		if !strings.Contains(received, msg) {
			t.Errorf("Response missing message: %s", msg)
		}
	}
}

func TestEventStreamResp_ChannelScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*Context) chan<- []byte
		verify      func(*testing.T, *httptest.ResponseRecorder, chan<- []byte) bool
		expectPanic bool
	}{
		{
			name: "Scenario 1: Send multiple data, channel auto-closed",
			setup: func(ctx *Context) chan<- []byte {
				return ctx.EventStreamResp()
			},
			verify: func(t *testing.T, resp *httptest.ResponseRecorder, eventCh chan<- []byte) bool {
				messages := []string{"msg1", "msg2", "msg3", "msg4"}

				for _, msg := range messages {
					eventCh <- []byte(msg)
				}

				time.Sleep(100 * time.Millisecond)

				received := resp.Body.String()
				for _, msg := range messages {
					if !strings.Contains(received, msg) {
						t.Errorf("Expected message %q not found in response", msg)
					}
				}
				var isPanic bool
				return isPanic
			},
			expectPanic: false,
		},
		{
			name: "Scenario 2: Send data then manually close channel",
			setup: func(ctx *Context) chan<- []byte {
				return ctx.EventStreamResp()
			},
			verify: func(t *testing.T, resp *httptest.ResponseRecorder, eventCh chan<- []byte) (isPanic bool) {
				eventCh <- []byte("first message")
				time.Sleep(10 * time.Millisecond)
				if !strings.Contains(resp.Body.String(), "first message") {
					t.Error("First message not received")
				}
				close(eventCh)
				func() {
					defer func() {
						if r := recover(); r != nil {
							isPanic = true
						}
					}()
					eventCh <- []byte("should panic")
				}()
				return
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/events", nil)
			resp := httptest.NewRecorder()

			ctxInstance := &Context{
				Request:        req,
				ResponseWriter: &Response{ResponseWriter: resp},
			}

			eventCh := tt.setup(ctxInstance)

			resPanicStatus := tt.verify(t, resp, eventCh)
			assert.Equal(t, tt.expectPanic, resPanicStatus)

		})
	}
}

func TestXsrfReset_01(t *testing.T) {
	r := &http.Request{}
	c := NewContext()
	c.Request = r
	c.ResponseWriter = &Response{}
	c.ResponseWriter.reset(httptest.NewRecorder())
	c.Output.Reset(c)
	c.Input.Reset(c)
	c.XSRFToken("key", 16)
	if c._xsrfToken == "" {
		t.FailNow()
	}
	token := c._xsrfToken
	c.Reset(&Response{ResponseWriter: httptest.NewRecorder()}, r)
	if c._xsrfToken != "" {
		t.FailNow()
	}
	c.XSRFToken("key", 16)
	if c._xsrfToken == "" {
		t.FailNow()
	}
	if token == c._xsrfToken {
		t.FailNow()
	}
}

func TestContext_Session(t *testing.T) {
	c := NewContext()
	if store, err := c.Session(); store != nil || err == nil {
		t.FailNow()
	}
}

func TestContext_Session1(t *testing.T) {
	c := Context{}
	if store, err := c.Session(); store != nil || err == nil {
		t.FailNow()
	}
}

func TestContext_Session2(t *testing.T) {
	c := NewContext()
	c.Input.CruSession = &session.MemSessionStore{}

	if store, err := c.Session(); store == nil || err != nil {
		t.FailNow()
	}
}

func TestSetCookie(t *testing.T) {
	type cookie struct {
		Name     string
		Value    string
		MaxAge   int64
		Path     string
		Domain   string
		Secure   bool
		HttpOnly bool
		SameSite string
	}
	type testItem struct {
		item cookie
		want string
	}
	cases := []struct {
		request string
		valueGp []testItem
	}{
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "Strict"}, "name=value; Max-Age=0; Path=/; SameSite=Strict"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "Lax"}, "name=value; Max-Age=0; Path=/; SameSite=Lax"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "None"}, "name=value; Max-Age=0; Path=/; SameSite=None"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, ""}, "name=value; Max-Age=0; Path=/"}}},
	}
	for _, c := range cases {
		r, _ := http.NewRequest("GET", c.request, nil)
		output := NewOutput()
		output.Context = NewContext()
		output.Context.Reset(httptest.NewRecorder(), r)
		for _, item := range c.valueGp {
			params := item.item
			others := []interface{}{params.MaxAge, params.Path, params.Domain, params.Secure, params.HttpOnly, params.SameSite}
			output.Context.SetCookie(params.Name, params.Value, others...)
			got := output.Context.ResponseWriter.Header().Get("Set-Cookie")
			if got != item.want {
				t.Fatalf("SetCookie error,should be:\n%v \ngot:\n%v", item.want, got)
			}
		}
	}
}
