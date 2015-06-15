// Copyright 2014 beego Author. All Rights Reserved.
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

package beego

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

type TLSServer struct {
	*http.Server
	Certificate                 tls.Certificate
	TLSMinVersion               uint16
	TLSMaxVersion               uint16
	TLSPreferServerCipherSuites bool
	TLSCiphers                  []uint16
}

//Creating TLS configuration information
func (srv *TLSServer) TLSConfigServer() (err error) {
	srv.TLSPreferServerCipherSuites = TLSPreferServerCipher
	srv.TLSMinVersion = defaultsupportedProtocols[TLSMinVersion]
	srv.TLSMaxVersion = defaultsupportedProtocols[TLSMaxVersion]

	for _, value := range TLSCiphers {
		srv.TLSCiphers = append(srv.TLSCiphers, supportedCiphersMap[value])
	}
	if len(srv.TLSCiphers) < 1 {
		srv.TLSCiphers = defaultsupportedCiphers
	}

	if HttpsCertFile != "" && HttpsKeyFile != "" {
		var err error
		srv.Certificate, err = tls.LoadX509KeyPair(HttpsCertFile, HttpsKeyFile)
		if err != nil {
			return err
		}
		return nil
	}
	HttpsCertContent, HttpsKeyContent := AppConfig.String("HttpsCertContent"), AppConfig.String("HttpsKeyContent")
	srv.Certificate, err = tls.X509KeyPair([]byte(HttpsCertContent), []byte(HttpsKeyContent))
	if err != nil {
		return err
	}
	srv.TLSConfig = &tls.Config{
		CipherSuites:             srv.TLSCiphers,
		PreferServerCipherSuites: srv.TLSPreferServerCipherSuites,
		MinVersion:               srv.TLSMinVersion,
		MaxVersion:               srv.TLSMaxVersion,
	}
	srv.TLSConfig.Certificates = make([]tls.Certificate, 1)
	srv.TLSConfig.Certificates[0] = srv.Certificate
	return nil
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for
// the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
func (srv *TLSServer) ListenAndServeTLS() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.TLSConfig == nil {
		return errors.New("Not Configured TLSConfig")
	}
	if len(srv.TLSConfig.Certificates) < 1 {
		return errors.New("SSL certificate information is not configured")
	}
	*config = *srv.TLSConfig
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return srv.Serve(tlsListener)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// List of supported cipher suites in descending order of preference.
// Ordering is very important! Getting the wrong order will break
// Note that TLS_FALLBACK_SCSV is not in this list since it is always
// added manually.
var defaultsupportedCiphers = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
}

// Map of supported ciphers, used only for parsing config.
//
// Note that, at time of writing, HTTP/2 blacklists 276 cipher suites,
// including all but two of the suites below (the two GCM suites).
// See https://http2.github.io/http2-spec/#BadCipherSuites
//
// TLS_FALLBACK_SCSV is not in this list because we manually ensure
// it is always added (even though it is not technically a cipher suite).
//
// This map, like any map, is NOT ORDERED. Do not range over this map.
var supportedCiphersMap = map[string]uint16{
	"ECDHE-RSA-AES128-GCM-SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-ECDSA-AES128-GCM-SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	"ECDHE-RSA-AES128-CBC-SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	"ECDHE-RSA-AES256-CBC-SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	"ECDHE-ECDSA-AES256-CBC-SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	"ECDHE-ECDSA-AES128-CBC-SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	"RSA-AES128-CBC-SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	"RSA-AES256-CBC-SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	"ECDHE-RSA-3DES-EDE-CBC-SHA":    tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	"RSA-3DES-EDE-CBC-SHA":          tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
}

// Map of supported protocols
// SSLv3 will be not supported in future release
// HTTP/2 only supports TLS 1.2 and higher
var defaultsupportedProtocols = map[string]uint16{
	"ssl3.0": tls.VersionSSL30,
	"tls1.0": tls.VersionTLS10,
	"tls1.1": tls.VersionTLS11,
	"tls1.2": tls.VersionTLS12,
}
