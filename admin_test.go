package beego

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/astaxie/beego/toolbox"
)

type SampleDatabaseCheck struct {
}

type SampleCacheCheck struct {
}

func (dc *SampleDatabaseCheck) Check() error {
	return nil
}

func (cc *SampleCacheCheck) Check() error {
	return errors.New("no cache detected")
}

func TestList_01(t *testing.T) {
	m := make(M)
	list("BConfig", BConfig, m)
	t.Log(m)
	om := oldMap()
	for k, v := range om {
		if fmt.Sprint(m[k]) != fmt.Sprint(v) {
			t.Log(k, "old-key", v, "new-key", m[k])
			t.FailNow()
		}
	}
}

func oldMap() M {
	m := make(M)
	m["BConfig.AppName"] = BConfig.AppName
	m["BConfig.RunMode"] = BConfig.RunMode
	m["BConfig.RouterCaseSensitive"] = BConfig.RouterCaseSensitive
	m["BConfig.ServerName"] = BConfig.ServerName
	m["BConfig.RecoverPanic"] = BConfig.RecoverPanic
	m["BConfig.CopyRequestBody"] = BConfig.CopyRequestBody
	m["BConfig.EnableGzip"] = BConfig.EnableGzip
	m["BConfig.MaxMemory"] = BConfig.MaxMemory
	m["BConfig.EnableErrorsShow"] = BConfig.EnableErrorsShow
	m["BConfig.Listen.Graceful"] = BConfig.Listen.Graceful
	m["BConfig.Listen.ServerTimeOut"] = BConfig.Listen.ServerTimeOut
	m["BConfig.Listen.ListenTCP4"] = BConfig.Listen.ListenTCP4
	m["BConfig.Listen.EnableHTTP"] = BConfig.Listen.EnableHTTP
	m["BConfig.Listen.HTTPAddr"] = BConfig.Listen.HTTPAddr
	m["BConfig.Listen.HTTPPort"] = BConfig.Listen.HTTPPort
	m["BConfig.Listen.EnableHTTPS"] = BConfig.Listen.EnableHTTPS
	m["BConfig.Listen.HTTPSAddr"] = BConfig.Listen.HTTPSAddr
	m["BConfig.Listen.HTTPSPort"] = BConfig.Listen.HTTPSPort
	m["BConfig.Listen.HTTPSCertFile"] = BConfig.Listen.HTTPSCertFile
	m["BConfig.Listen.HTTPSKeyFile"] = BConfig.Listen.HTTPSKeyFile
	m["BConfig.Listen.EnableAdmin"] = BConfig.Listen.EnableAdmin
	m["BConfig.Listen.AdminAddr"] = BConfig.Listen.AdminAddr
	m["BConfig.Listen.AdminPort"] = BConfig.Listen.AdminPort
	m["BConfig.Listen.EnableFcgi"] = BConfig.Listen.EnableFcgi
	m["BConfig.Listen.EnableStdIo"] = BConfig.Listen.EnableStdIo
	m["BConfig.WebConfig.AutoRender"] = BConfig.WebConfig.AutoRender
	m["BConfig.WebConfig.EnableDocs"] = BConfig.WebConfig.EnableDocs
	m["BConfig.WebConfig.FlashName"] = BConfig.WebConfig.FlashName
	m["BConfig.WebConfig.FlashSeparator"] = BConfig.WebConfig.FlashSeparator
	m["BConfig.WebConfig.DirectoryIndex"] = BConfig.WebConfig.DirectoryIndex
	m["BConfig.WebConfig.StaticDir"] = BConfig.WebConfig.StaticDir
	m["BConfig.WebConfig.StaticExtensionsToGzip"] = BConfig.WebConfig.StaticExtensionsToGzip
	m["BConfig.WebConfig.StaticCacheFileSize"] = BConfig.WebConfig.StaticCacheFileSize
	m["BConfig.WebConfig.StaticCacheFileNum"] = BConfig.WebConfig.StaticCacheFileNum
	m["BConfig.WebConfig.TemplateLeft"] = BConfig.WebConfig.TemplateLeft
	m["BConfig.WebConfig.TemplateRight"] = BConfig.WebConfig.TemplateRight
	m["BConfig.WebConfig.ViewsPath"] = BConfig.WebConfig.ViewsPath
	m["BConfig.WebConfig.EnableXSRF"] = BConfig.WebConfig.EnableXSRF
	m["BConfig.WebConfig.XSRFExpire"] = BConfig.WebConfig.XSRFExpire
	m["BConfig.WebConfig.Session.SessionOn"] = BConfig.WebConfig.Session.SessionOn
	m["BConfig.WebConfig.Session.SessionProvider"] = BConfig.WebConfig.Session.SessionProvider
	m["BConfig.WebConfig.Session.SessionName"] = BConfig.WebConfig.Session.SessionName
	m["BConfig.WebConfig.Session.SessionGCMaxLifetime"] = BConfig.WebConfig.Session.SessionGCMaxLifetime
	m["BConfig.WebConfig.Session.SessionProviderConfig"] = BConfig.WebConfig.Session.SessionProviderConfig
	m["BConfig.WebConfig.Session.SessionCookieLifeTime"] = BConfig.WebConfig.Session.SessionCookieLifeTime
	m["BConfig.WebConfig.Session.SessionAutoSetCookie"] = BConfig.WebConfig.Session.SessionAutoSetCookie
	m["BConfig.WebConfig.Session.SessionDomain"] = BConfig.WebConfig.Session.SessionDomain
	m["BConfig.WebConfig.Session.SessionDisableHTTPOnly"] = BConfig.WebConfig.Session.SessionDisableHTTPOnly
	m["BConfig.Log.AccessLogs"] = BConfig.Log.AccessLogs
	m["BConfig.Log.EnableStaticLogs"] = BConfig.Log.EnableStaticLogs
	m["BConfig.Log.AccessLogsFormat"] = BConfig.Log.AccessLogsFormat
	m["BConfig.Log.FileLineNum"] = BConfig.Log.FileLineNum
	m["BConfig.Log.Outputs"] = BConfig.Log.Outputs
	return m
}

func TestWriteJSON(t *testing.T) {
	t.Log("Testing the adding of JSON to the response")

	w := httptest.NewRecorder()
	originalBody := []int{1, 2, 3}

	res, _ := json.Marshal(originalBody)

	writeJSON(w, res)

	decodedBody := []int{}
	err := json.NewDecoder(w.Body).Decode(&decodedBody)

	if err != nil {
		t.Fatal("Could not decode response body into slice.")
	}

	for i := range decodedBody {
		if decodedBody[i] != originalBody[i] {
			t.Fatalf("Expected %d but got %d in decoded body slice", originalBody[i], decodedBody[i])
		}
	}
}

func TestHealthCheckHandlerDefault(t *testing.T) {
	endpointPath := "/healthcheck"

	toolbox.AddHealthCheck("database", &SampleDatabaseCheck{})
	toolbox.AddHealthCheck("cache", &SampleCacheCheck{})

	req, err := http.NewRequest("GET", endpointPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(healthcheck)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "database") {
		t.Errorf("Expected 'database' in generated template.")
	}

}

func TestBuildHealthCheckResponseList(t *testing.T) {
	healthCheckResults := [][]string{
		{
			"error",
			"Database",
			"Error occurred whie starting the db",
		},
		{
			"success",
			"Cache",
			"Cache started successfully",
		},
	}

	responseList := buildHealthCheckResponseList(&healthCheckResults)

	if len(responseList) != len(healthCheckResults) {
		t.Errorf("invalid response map length: got %d want %d",
			len(responseList), len(healthCheckResults))
	}

	responseFields := []string{"name", "message", "status"}

	for _, response := range responseList {
		for _, field := range responseFields {
			_, ok := response[field]
			if !ok {
				t.Errorf("expected %s to be in the response %v", field, response)
			}
		}

	}

}

func TestHealthCheckHandlerReturnsJSON(t *testing.T) {

	toolbox.AddHealthCheck("database", &SampleDatabaseCheck{})
	toolbox.AddHealthCheck("cache", &SampleCacheCheck{})

	req, err := http.NewRequest("GET", "/healthcheck?json=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(healthcheck)

	handler.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	decodedResponseBody := []map[string]interface{}{}
	expectedResponseBody := []map[string]interface{}{}

	expectedJSONString := []byte(`
		[
			{
				"message":"database",
				"name":"success",
				"status":"OK"
			},
			{
				"message":"cache",
				"name":"error",
				"status":"no cache detected"
			}
		]
	`)

	json.Unmarshal(expectedJSONString, &expectedResponseBody)

	json.Unmarshal(w.Body.Bytes(), &decodedResponseBody)

	if len(expectedResponseBody) != len(decodedResponseBody) {
		t.Errorf("invalid response map length: got %d want %d",
			len(decodedResponseBody), len(expectedResponseBody))
	}
	assert.Equal(t, len(expectedResponseBody), len(decodedResponseBody))
	assert.Equal(t, 2, len(decodedResponseBody))

	var database, cache map[string]interface{}
	if decodedResponseBody[0]["message"] == "database" {
		database = decodedResponseBody[0]
		cache = decodedResponseBody[1]
	} else {
		database = decodedResponseBody[1]
		cache = decodedResponseBody[0]
	}

	assert.Equal(t, expectedResponseBody[0], database)
	assert.Equal(t, expectedResponseBody[1], cache)
}
