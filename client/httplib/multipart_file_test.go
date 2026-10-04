// Copyright 2026 beego Author. All Rights Reserved.
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
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// multipartFileServer serves upload requests over TCP with bounded deadlines.
func multipartFileServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewUnstartedServer(handler)
	server.Config.ReadTimeout = 5 * time.Second
	server.Config.WriteTimeout = 5 * time.Second
	server.Config.ConnContext = func(ctx context.Context, conn net.Conn) context.Context {
		if _, ok := conn.(*net.TCPConn); !ok {
			t.Errorf("upload connection is %T, want native TCP", conn)
		}
		return ctx
	}
	server.Start()
	t.Cleanup(server.Close)
	return server
}

// TestPostFileUnreadable checks file errors for every multipart request method.
func TestPostFileUnreadable(t *testing.T) {
	for _, method := range []string{"POST", "PUT", "PATCH", "DELETE"} {
		for _, failure := range []string{"missing", "directory"} {
			t.Run(method+"/"+failure, func(t *testing.T) {
				assertPostFileUnreadable(t, method, failure)
			})
		}
	}
}

// assertPostFileUnreadable verifies the original filesystem error reaches Response.
func assertPostFileUnreadable(t *testing.T, method, failure string) {
	t.Helper()
	filename := t.TempDir()
	if failure == "missing" {
		filename = filepath.Join(filename, "absent.bin")
	}
	server := multipartFileServer(t, func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(io.Discard, r.Body); err == nil {
			w.WriteHeader(http.StatusOK)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	transport := &http.Transport{}
	defer transport.CloseIdleConnections()
	request := NewBeegoRequestWithCtx(ctx, server.URL, method).
		SetTransport(transport).PostFile("document", filename)
	response, err := request.Response()
	if response != nil {
		response.Body.Close()
	}
	if err == nil {
		t.Fatalf("upload of %s returned success, want file error", filename)
	}
	var pathError *os.PathError
	if !errors.As(err, &pathError) || pathError.Path != filename {
		t.Fatalf("upload error = %v, want original PathError for %s", err, filename)
	}
	if failure == "missing" && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("upload error = %v, want os.ErrNotExist", err)
	}
}

// multipartFileUpload records the server-observed method, fields and file data.
type multipartFileUpload struct {
	method string
	fields url.Values
	files  map[string][]byte
	names  map[string]string
}

// TestPostFileMultipartSuccess checks valid files and form fields for every method.
func TestPostFileMultipartSuccess(t *testing.T) {
	for _, method := range []string{"POST", "PUT", "PATCH", "DELETE"} {
		t.Run(method, func(t *testing.T) {
			assertMultipartFileUpload(t, method)
		})
	}
}

// multipartUploadHandler reads the actual multipart files and fields on the server.
func multipartUploadHandler(received chan<- multipartFileUpload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.MultipartForm.RemoveAll()
		upload := multipartFileUpload{r.Method, r.MultipartForm.Value,
			make(map[string][]byte), make(map[string]string)}
		for name, headers := range r.MultipartForm.File {
			if len(headers) != 1 {
				http.Error(w, "unexpected file count", http.StatusBadRequest)
				return
			}
			file, err := headers[0].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			data, readErr := io.ReadAll(file)
			closeErr := file.Close()
			if readErr != nil || closeErr != nil {
				http.Error(w, "could not read uploaded file", http.StatusBadRequest)
				return
			}
			upload.files[name] = data
			upload.names[name] = headers[0].Filename
		}
		received <- upload
		w.WriteHeader(http.StatusOK)
	}
}

// assertMultipartFileUpload compares uploaded bytes, filenames and repeated fields.
func assertMultipartFileUpload(t *testing.T, method string) {
	t.Helper()
	directory := t.TempDir()
	binaryName := filepath.Join(directory, "payload.bin")
	emptyName := filepath.Join(directory, "empty.txt")
	binaryData := bytes.Repeat([]byte{0, 255, '\n', 128}, 1025)
	if err := os.WriteFile(binaryName, binaryData, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(emptyName, nil, 0600); err != nil {
		t.Fatal(err)
	}
	received := make(chan multipartFileUpload, 1)
	server := multipartFileServer(t, multipartUploadHandler(received))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	transport := &http.Transport{}
	defer transport.CloseIdleConnections()
	request := NewBeegoRequestWithCtx(ctx, server.URL, method).SetTransport(transport).
		PostFile("binary", binaryName).PostFile("empty-file", emptyName).
		Param("note", "中文 value").Param("note", "second").Param("empty", "")
	response, err := request.Response()
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("upload status = %d, want 200", response.StatusCode)
	}
	upload := <-received
	assertReceivedMultipartFile(t, upload, method, binaryData)
}

// assertReceivedMultipartFile checks the server-observed payload against the input.
func assertReceivedMultipartFile(t *testing.T, upload multipartFileUpload, method string, binaryData []byte) {
	t.Helper()
	if upload.method != method || len(upload.files) != 2 {
		t.Fatalf("upload method/files = %s/%d", upload.method, len(upload.files))
	}
	if !bytes.Equal(upload.files["binary"], binaryData) || len(upload.files["empty-file"]) != 0 {
		t.Fatal("multipart upload changed binary or empty file content")
	}
	wantNames := map[string]string{"binary": "payload.bin", "empty-file": "empty.txt"}
	if !reflect.DeepEqual(upload.names, wantNames) {
		t.Fatalf("file names = %v, want %v", upload.names, wantNames)
	}
	wantFields := url.Values{"note": {"中文 value", "second"}, "empty": {""}}
	if !reflect.DeepEqual(upload.fields, wantFields) {
		t.Fatalf("fields = %v, want %v", upload.fields, wantFields)
	}
}

// TestMultipartBodyFileError checks that reading the multipart pipe preserves errors.
func TestMultipartBodyFileError(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "absent.bin")
	request := Post("http://localhost/").PostFile("document", filename)
	request.handleFiles()
	defer request.GetRequest().Body.Close()
	_, err := io.ReadAll(request.GetRequest().Body)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("multipart stream error = %v, want os.ErrNotExist", err)
	}
}
