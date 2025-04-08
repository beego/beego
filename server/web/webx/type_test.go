package webx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAllWrapperTestCase(t *testing.T) {

	type MyStruct struct {
		Foo string `form:"foo" json:"foo"`
	}

	myStruct := MyStruct{Foo: "bar"}
	webFunc := func(_ *context.Context, s1 MyStruct) (any, error) {
		return s1, nil
	}
	testCases := []struct {
		name         string
		expectedCode int
		expectedRes  func() string
		reqBody      func() io.Reader
		contentType  string
		bizProvider  func() web.HandleFunc
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
			bizProvider: func() web.HandleFunc {
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
			bizProvider: func() web.HandleFunc {
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
			bizProvider: func() web.HandleFunc {
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
			bizProvider: func() web.HandleFunc {
				return internalWrapper(webFunc, func(_ *context.Context, _ *MyStruct) error {
					return errors.New("paras entity error")
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
			bizProvider: func() web.HandleFunc {
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
			app := web.NewHttpSever()
			// 2. set copy request body
			app.Cfg.CopyRequestBody = true

			// 3. register route
			path := "/api/data"
			app.Post(path, tc.bizProvider())

			// 4. create request
			req := httptest.NewRequest("POST", path, tc.reqBody())
			req.Header.Set("Content-Type", tc.contentType)
			req.Header.Set("Accept", "*/*")
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
