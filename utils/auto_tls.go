package utils

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
)

func AutoTLS(dirCache string, hosts ...string) tls.Config {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Cache:      autocert.DirCache(dirCache),
	}
	return tls.Config{GetCertificate: m.GetCertificate}
}
