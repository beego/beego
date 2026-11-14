// Copyright 2026
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

package web

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

// TestRunHTTPAndHTTPSSeparateListeners is a regression test for
// https://github.com/beego/beego/issues/5930: in normal (non-graceful)
// mode the HTTP and HTTPS listeners shared the same *http.Server, so the
// two startup goroutines raced on Server.Addr and a listener could bind
// the other protocol's port. Each listener must serve its own protocol
// on its own port. Run with -race: the shared-server code reports a data
// race on http.Server.Addr.
func TestRunHTTPAndHTTPSSeparateListeners(t *testing.T) {
	// Run executes the start hooks, and registerTemplate locks the view
	// paths globally. Restore the flag so tests running after this one
	// can still add view paths.
	t.Cleanup(func() { beeViewPathTemplateLocked = false })

	certFile, keyFile := writeSelfSignedCert(t)

	httpPort := freePort(t)
	httpsPort := freePort(t)

	cfg := *BConfig
	cfg.Listen.Graceful = false
	cfg.Listen.ListenTCP4 = false
	cfg.Listen.EnableHTTP = true
	cfg.Listen.EnableHTTPS = true
	cfg.Listen.EnableMutualHTTPS = false
	cfg.Listen.AutoTLS = false
	cfg.Listen.EnableAdmin = false
	cfg.Listen.EnableFcgi = false
	cfg.Listen.HTTPAddr = "127.0.0.1"
	cfg.Listen.HTTPPort = httpPort
	cfg.Listen.HTTPSAddr = "127.0.0.1"
	cfg.Listen.HTTPSPort = httpsPort
	cfg.Listen.HTTPSCertFile = certFile
	cfg.Listen.HTTPSKeyFile = keyFile

	app := NewHttpServerWithCfg(&cfg)
	app.Get("/ping", func(ctx *beecontext.Context) {
		_ = ctx.Output.Body([]byte("pong"))
	})

	go app.Run(fmt.Sprintf("127.0.0.1:%d", httpPort))

	httpURL := fmt.Sprintf("http://127.0.0.1:%d/ping", httpPort)
	httpsURL := fmt.Sprintf("https://127.0.0.1:%d/ping", httpsPort)
	httpsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // self-signed test cert
		},
		Timeout: 2 * time.Second,
	}
	httpClient := &http.Client{Timeout: 2 * time.Second}

	deadline := time.Now().Add(20 * time.Second)
	for {
		httpOK := fetchBody(t, httpClient, httpURL) == "pong"
		httpsOK := fetchBody(t, httpsClient, httpsURL) == "pong"
		if httpOK && httpsOK {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("listeners did not come up: http ok=%v https ok=%v", httpOK, httpsOK)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// The HTTPS port must speak TLS, not plaintext HTTP served by the
	// HTTP listener that stole the port.
	resp, err := httpClient.Get(httpsURL)
	if err == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("plaintext HTTP request to the HTTPS port unexpectedly succeeded")
		}
	}
}

func fetchBody(t *testing.T, client *http.Client, url string) string {
	t.Helper()
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate a free port: %v", err)
	}
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port
}

func writeSelfSignedCert(t *testing.T) (certFile, keyFile string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:     []string{"localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("failed to create test certificate: %v", err)
	}
	dir := t.TempDir()
	certFile = filepath.Join(dir, "cert.pem")
	keyFile = filepath.Join(dir, "key.pem")
	certOut, err := os.Create(certFile)
	if err != nil {
		t.Fatalf("failed to create test cert file: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		t.Fatalf("failed to write test cert: %v", err)
	}
	_ = certOut.Close()
	keyDer, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal test key: %v", err)
	}
	keyOut, err := os.Create(keyFile)
	if err != nil {
		t.Fatalf("failed to create test key file: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDer}); err != nil {
		t.Fatalf("failed to write test key: %v", err)
	}
	_ = keyOut.Close()
	return certFile, keyFile
}
