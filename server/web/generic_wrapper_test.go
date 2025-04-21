package web

import (
	"bytes"
	ctx0 "context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// userRequest is a struct that represents the request parameters.
type userRequest struct {
	Name string `json:"name" form:"name"`
	Age  int    `json:"age"`
}

// addUser is a sample business logic function that takes a context and userRequest as parameters.
func addUser(_ *context.Context, params userRequest) (any, error) {
	if params.Name == "" {
		return nil, errors.New("name can't be null")
	}
	return []any{params.Name, params.Age}, nil
}

// TestWrapperFromJsonExample test the WrapperFromJson function.
func TestWrapperFromJsonExample(t *testing.T) {

	app := NewHttpSever()
	app.Cfg.CopyRequestBody = true
	path := "/api/data"
	// to use wrapper
	app.Post(path, Wrapper(addUser))

	reader := strings.NewReader(`{"name": "rose", "age": 17}`)

	req := httptest.NewRequest("POST", path, reader)
	req.Header.Set("Content-Type", context.ApplicationJSON)
	req.Header.Set("Accept", "*/*")

	w := httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	if w.Code == http.StatusOK {
		marshal, _ := json.Marshal([]any{"rose", 17})
		assert.Equal(t, string(marshal), w.Body.String())
	}
}

// TestWrapperFromFormExample test the WrapperFromForm function.
func TestWrapperFromFormExample(t *testing.T) {

	app := NewHttpSever()
	app.Cfg.CopyRequestBody = true
	path := "/api/data"
	// to use wrapper
	app.Post(path, Wrapper(addUser))

	formData := url.Values{}
	formData.Set("name", "jack")

	req := httptest.NewRequest("POST", path, strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", context.ApplicationForm)
	req.Header.Set("Accept", "*/*")

	w := httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	if w.Code == http.StatusOK {
		marshal, _ := json.Marshal([]any{"jack", 0})
		assert.Equal(t, string(marshal), w.Body.String())
	}
}

// TestWrapperExample test the Wrapper function.
func TestWrapperExample(t *testing.T) {

	app := NewHttpSever()
	app.Cfg.CopyRequestBody = true
	path := "/api/data"
	// to use wrapper
	app.Post(path, Wrapper(addUser))

	request := userRequest{
		Name: "tom",
		Age:  18,
	}
	marshal, _ := xml.Marshal(request)

	req := httptest.NewRequest("POST", path, bytes.NewBuffer(marshal))
	req.Header.Set("Content-Type", context.ApplicationXML)
	req.Header.Set("Accept", "*/*")

	w := httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	if w.Code == http.StatusOK {
		marshal, _ := json.Marshal([]any{"tom", 18})
		assert.Equal(t, string(marshal), w.Body.String())
	}
}

func TestAllWrapperTestCase(t *testing.T) {

	type MyStruct struct {
		Foo string `form:"foo" json:"foo"`
	}

	myStruct := MyStruct{Foo: "bar"}

	webFunc := func(_ *context.Context, s1 MyStruct) (any, error) {
		return s1, nil
	}

	testCases := []struct {
		name              string
		expectedCode      int
		expectedRes       func() string
		reqBody           func() io.Reader
		reqHeader         map[string]string
		useDefaultSession bool
		contentType       string
		bizProvider       func() HandleFunc
	}{
		{
			name:         "Test post json requestBody",
			expectedCode: http.StatusOK,
			expectedRes: func() string {
				marshal, _ := json.Marshal(myStruct)
				return string(marshal)
			},
			reqBody: func() io.Reader {
				return strings.NewReader(`{"foo": "bar"}`)
			},
			contentType: context.ApplicationJSON,
			bizProvider: func() HandleFunc {
				return WrapperFromJson(webFunc)
			},
		},
		{
			name:         "Test post form requestBody",
			expectedCode: http.StatusOK,
			expectedRes: func() string {
				marshal, _ := json.Marshal(myStruct)
				return string(marshal)
			},
			reqBody: func() io.Reader {
				formData := url.Values{}
				formData.Set("foo", "bar")
				return strings.NewReader(formData.Encode())
			},
			contentType: context.ApplicationForm,
			bizProvider: func() HandleFunc {
				return WrapperFromForm(webFunc)
			},
		},
		{
			name:         "Test base binging",
			expectedCode: http.StatusOK,
			expectedRes: func() string {
				marshal, _ := json.Marshal(myStruct)
				return string(marshal)
			},
			reqBody: func() io.Reader {
				marshal, _ := xml.Marshal(myStruct)
				return bytes.NewBuffer(marshal)
			},
			contentType: context.ApplicationXML,
			bizProvider: func() HandleFunc {
				return Wrapper(webFunc)
			},
		},
		{
			name:         "Test unWrapper error",
			expectedCode: http.StatusBadRequest,
			reqBody: func() io.Reader {
				formData := url.Values{}
				formData.Set("foo", "bar")
				return strings.NewReader(formData.Encode())
			},
			contentType: context.ApplicationForm,
			bizProvider: func() HandleFunc {
				return internalWrapper(webFunc, func(ctx *context.Context) (params MyStruct, err error) {
					err = errors.New("paras entity error")
					return
				})
			},
		},
		{
			name:         "Test biz error",
			expectedCode: http.StatusInternalServerError,
			reqBody: func() io.Reader {
				formData := url.Values{}
				formData.Set("foo", "bar")
				return strings.NewReader(formData.Encode())
			},
			contentType: context.ApplicationForm,
			bizProvider: func() HandleFunc {
				testFunc := func(_ *context.Context, _ MyStruct) (any, error) {
					return nil, errors.New("biz error")
				}
				return WrapperFromForm(testFunc)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// 1. init web server
			app := NewHttpSever()
			// 2. set copy request body
			app.Cfg.CopyRequestBody = true
			// need to config session before register route
			if tc.useDefaultSession {
				c := defaultSessionOption(app)
				// clear session config
				defer c()
			}
			// 3. register route
			path := "/api/data"
			app.Post(path, tc.bizProvider())

			// 4. create request
			req := httptest.NewRequest("POST", path, tc.reqBody())
			req.Header.Set("Content-Type", tc.contentType)
			req.Header.Set("Accept", "*/*")

			for key, value := range tc.reqHeader {
				req.Header.Set(key, value)
			}

			// 5. create ResponseRecorder
			w := httptest.NewRecorder()

			// 6. process request
			app.Handlers.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)
			if w.Code == http.StatusOK {
				assert.Equal(t, tc.expectedRes(), w.Body.String())
			}
		})
	}

}

type userInfo struct {
	ID       int
	Username string
	Role     string
}

const sessionKey = "user_info"

var defaultUser = userInfo{
	ID:       0,
	Username: "guest",
	Role:     "guest",
}

func defaultSessionOption(app *HttpServer) (cancel func()) {

	config := `{"cookieName":"gosessionid","enableSetCookie":false,"gclifetime":3600,"ProviderConfig":"{\"cookieName\":\"gosessionid\",\"securityKey\":\"beegocookiehashkey\"}"}`
	conf := new(session.ManagerConfig)
	_ = json.Unmarshal([]byte(config), conf)
	GlobalSessions, _ = session.NewManager("cookie", conf)

	app.Cfg.WebConfig.Session.SessionOn = true
	app.Cfg.WebConfig.Session.SessionProvider = "memory"
	app.Cfg.WebConfig.Session.SessionName = "beegoSessionId"
	app.Cfg.WebConfig.Session.SessionGCMaxLifetime = 3600

	app.InsertFilter("*", BeforeExec, func(ctx *context.Context) {
		if ctx.Input.Session(sessionKey) == nil {
			timeout, c := ctx0.WithTimeout(ctx0.Background(), time.Minute*10)
			defer c()
			_ = ctx.Input.CruSession.Set(timeout, "user_info", &defaultUser)
		}
	})

	return func() {
		GlobalSessions = nil
		app.Cfg.WebConfig.Session.SessionOn = false
	}
}
