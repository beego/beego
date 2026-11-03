// Copyright 2026 beego Author. All Rights Reserved.
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

package web

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

func freeLoopbackPort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func writeSelfSignedCert(t *testing.T, dir string) (certFile, keyFile string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "127.0.0.1"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)

	certFile = filepath.Join(dir, "cert.pem")
	keyFile = filepath.Join(dir, "key.pem")
	require.NoError(t, os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600))
	require.NoError(t, os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600))
	return certFile, keyFile
}

// Normal mode must serve HTTP and HTTPS on their own ports. Run it with
// -race: the two listeners used to write the same http.Server.Addr.
func TestRunNormalModeHTTPAndHTTPS(t *testing.T) {
	certFile, keyFile := writeSelfSignedCert(t, t.TempDir())
	httpPort := freeLoopbackPort(t)
	httpsPort := freeLoopbackPort(t)

	cfg := *BConfig
	cfg.Listen.Graceful = false
	cfg.Listen.ListenTCP4 = false
	cfg.Listen.EnableAdmin = false
	cfg.Listen.EnableHTTP = true
	cfg.Listen.HTTPAddr = "127.0.0.1"
	cfg.Listen.HTTPPort = httpPort
	cfg.Listen.EnableHTTPS = true
	cfg.Listen.HTTPSAddr = "127.0.0.1"
	cfg.Listen.HTTPSPort = httpsPort
	cfg.Listen.HTTPSCertFile = certFile
	cfg.Listen.HTTPSKeyFile = keyFile

	app := NewHttpServerWithCfg(&cfg)
	app.Handlers.Get("/ping", func(ctx *beecontext.Context) {
		_ = ctx.Output.Body([]byte("pong"))
	})

	app.Handlers.Init()
	app.Server.Handler = app.Handlers

	done := make(chan struct{})
	go func() {
		defer close(done)
		app.runNormalMode(fmt.Sprintf("127.0.0.1:%d", httpPort))
	}()

	// Trust the test certificate instead of skipping verification.
	certPEM, err := os.ReadFile(certFile)
	require.NoError(t, err)
	roots := x509.NewCertPool()
	require.True(t, roots.AppendCertsFromPEM(certPEM))
	client := &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			RootCAs:    roots,
			MinVersion: tls.VersionTLS12,
		}},
	}
	get := func(url string) string {
		var lastErr error
		for i := 0; i < 100; i++ {
			resp, err := client.Get(url)
			if err == nil {
				defer resp.Body.Close()
				buf := make([]byte, 16)
				n, _ := resp.Body.Read(buf)
				return string(buf[:n])
			}
			lastErr = err
			time.Sleep(50 * time.Millisecond)
		}
		t.Fatalf("GET %s: %v", url, lastErr)
		return ""
	}

	assert.Equal(t, "pong", get(fmt.Sprintf("http://127.0.0.1:%d/ping", httpPort)))
	assert.Equal(t, "pong", get(fmt.Sprintf("https://127.0.0.1:%d/ping", httpsPort)))

	// Closing the HTTP server makes runNormalMode return.
	require.NoError(t, app.Server.Shutdown(context.Background()))
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("runNormalMode did not return after the HTTP server was shut down")
	}
}

func TestCloneHTTPServer(t *testing.T) {
	src := &http.Server{
		Addr:              ":1",
		Handler:           http.NotFoundHandler(),
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       4 * time.Second,
		MaxHeaderBytes:    5,
		TLSConfig:         &tls.Config{ServerName: "x", MinVersion: tls.VersionTLS12},
	}
	dst := cloneHTTPServer(src)
	assert.Equal(t, src.Addr, dst.Addr)
	assert.NotNil(t, dst.Handler)
	assert.Equal(t, src.ReadTimeout, dst.ReadTimeout)
	assert.Equal(t, src.ReadHeaderTimeout, dst.ReadHeaderTimeout)
	assert.Equal(t, src.WriteTimeout, dst.WriteTimeout)
	assert.Equal(t, src.IdleTimeout, dst.IdleTimeout)
	assert.Equal(t, src.MaxHeaderBytes, dst.MaxHeaderBytes)
	assert.Equal(t, "x", dst.TLSConfig.ServerName)

	dst.Addr = ":2"
	dst.TLSConfig.ServerName = "y"
	assert.Equal(t, ":1", src.Addr)
	assert.Equal(t, "x", src.TLSConfig.ServerName)

	assert.Nil(t, cloneHTTPServer(&http.Server{}).TLSConfig)
}
