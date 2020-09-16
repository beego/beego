package httplib

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"
	"testing"
)

const clientCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIICVjCCAT4CAQAwETEPMA0GA1UEAwwGY2xpZW50MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAyxZb08Xd1IRAlIMfxBDqWCKTuvoEIfhVWp2zmREZRssx
xlkOLl9wBdxGtbB4Yii5xaqldIogXIyes2k8RKNeVycONv8whSBOjoxOcl1HIyUR
kQXmuyDeOQC6v+iAhs7a8LzubZl5Sq0/m1XcbANaHTa792IGXU9jx2u0xxip+pqP
lu3PdyQ3xpAeMGTTpIIxQu7ibs2l8Bq7aVrUlE3MLtZLRHxi+hiCJRkl9+M/NuaZ
wVYtJKHbnNlH3b54s+vKfLILzFC3CdAci4JWR5x2M+3x0JpMHUmwm0Dg8MgiAiLd
t8kTwTSGx4ZFhqY0AbU/7BAY8xtn6lM2lVsVthpMlQIDAQABoAAwDQYJKoZIhvcN
AQELBQADggEBALRAg+PO2kWphToTZnWx8lZ5vO10LXrOXAhi+KWUf678HIGaVaUG
1Z5UoLAPW5PHVqWDuFzjeEqSAECbCxdLkPpptz3QtUn/fq4LxWqNQTDZTcDCC35R
iTVRRTVC4t0l3M0U1hSrlp7VIVZtaGJABI6124ANm/7PtGqU7+kMryXxvaY1WRA+
H9bVJoyvfUWcAzpV+FFQspGe4t2xu2ot2a1jW1SHFXc8o4XxxD+64BBn6o//X6m0
lrmmZ+JC2QVmBAv0rRAYh4CxyOFQJnAW2wnmswXo8z0ZY66p/aVhTUCS25q9VUFD
Vnh1DhXn2bug786cqTIa99bReeoD7KzzHCg=
-----END CERTIFICATE REQUEST-----`
const clientCrt = `-----BEGIN CERTIFICATE-----
MIIDtjCCAZ6gAwIBAgIBAjANBgkqhkiG9w0BAQUFADATMREwDwYDVQQDDAh3YW5n
LmNvbTAeFw0yMDA5MTYwODI3NDNaFw0yMTA5MTYwODI3NDNaMBExDzANBgNVBAMM
BmNsaWVudDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMsWW9PF3dSE
QJSDH8QQ6lgik7r6BCH4VVqds5kRGUbLMcZZDi5fcAXcRrWweGIoucWqpXSKIFyM
nrNpPESjXlcnDjb/MIUgTo6MTnJdRyMlEZEF5rsg3jkAur/ogIbO2vC87m2ZeUqt
P5tV3GwDWh02u/diBl1PY8drtMcYqfqaj5btz3ckN8aQHjBk06SCMULu4m7NpfAa
u2la1JRNzC7WS0R8YvoYgiUZJffjPzbmmcFWLSSh25zZR92+eLPrynyyC8xQtwnQ
HIuCVkecdjPt8dCaTB1JsJtA4PDIIgIi3bfJE8E0hseGRYamNAG1P+wQGPMbZ+pT
NpVbFbYaTJUCAwEAAaMXMBUwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDQYJKoZIhvcN
AQEFBQADggIBAG+YR/M3RG8ls3hhoVCRkkhQ6ieJWt+zDCcao4lJbtzzalbDf9H8
TlUaTnrDrnBYVOPwdKtFrSricHeCXqv+hiilEMKGwv/77ScMDUNN0OfeHfXeqTcW
Lzg3NLjK1IX7dqBEiIvSEfd7VVY1/EQhVA7n+S0wrvDrSZvmOcxc+z/vV+y2a3en
etnJZ4Vy2eUD9Gt7hJO5nbzrN7WJQgb/WEGvuZlPHNa9Bx1kLkYJDMyCG2r5rr6n
fyFuRpZe/qrW201F5ga0DI1mG7jocIvqv6Ie5g8HabHR477Ub/OdgvyzX17fqO5n
YLYIK8HblJ2N51GQa4e11vWUo9b7ogFbanhJdytRKpYLj9aq02LmPhyVIDbwU/Bf
4bLac0Bg+7xpU70kUoy+gQtN3qtMSO+ZOXRmgfjbEESd0R8V0vfVeQtBqoIURvbp
6nDtXyLGofdzEWTL9Ab3HJL9nYQRZ898DEI5xBt2GDp3wS3vRIQEEILwzwX/nqrw
pwFqgdCLljp/1l8p2hR9r5U1edvUNjTL+l8zxgPyvhnDdI+cnEXTlV935esOfqdq
zht1HIu8cBLEf1bI2DJOc1CbhmdyH+2/dN/WV4uQp0kKQ7DrUZ8Lo053Gsl+xevo
xFs+72FYfZ4wMArSkbR3kEic3lNxwCG0JxBax66PrjGyqX1SgucqMefs
-----END CERTIFICATE-----
`
const caCrt = `-----BEGIN CERTIFICATE-----
MIIEojCCAooCCQC078JpAjeMtjANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAh3
YW5nLmNvbTAeFw0yMDA5MTYwODI3MDVaFw0yMTA5MTYwODI3MDVaMBMxETAPBgNV
BAMMCHdhbmcuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA3wN2
66llq/732RVSmDO7CW12E+MkIgNfvkJmLfB+DKxIZKYzZNMGnGxoBvjndL1PqJKt
IBb3hihtiLBLH11JfBGFKQfCntzJ3ynCOVlvn5a0aWa9t0rteOFGE8lFRfcUf7hj
UkJG4Zg7owmvQq9wWY7jgga99hc0Q3g2FUXThVhITsZeQiQy4kpceDgo8wTDiBlq
ltWtfm1RRsCvHRir9JAsa6IJlvTH/hkoSMyMsggP3a1LU2Zlo2T41CvK8HhEwK03
WmBhEBdkzzMTrxgZxcR2DoVBtked7mJuEHYG4U/LfQQJvW07Rh85dm6uhhUb381A
sRLcN5dpIfvGy1zKf5kNAZT3qKRsmjmCyFwyzkma97622++WegPponZ/xJLO0Ap4
hEkHdw+5YCOvnSKEi63/C6en1WS/OVLJtk+sw3CKL6T0sEWO8gwQti0mx4izGowL
xPvgyeltSsv0PgZpyRwybV9KfHz9Sa4rmKBV6071MTnFJ2b41aOziCFrYaMF45kz
TBqAMQRQIt/HH+Aimcog+XJtNDKOdOWEHpfw3/2ZRnJKXjsXeHdJQ7r/unQeGuA+
h+jKt73nHRVFolqLsQCjsZkyBJ1h3miL/0hRqhr/hEKeHvrsXGMlcom9rnRhAc7r
xyBPIgGkMXahacuUz1y2rb1S/t1SFBPJ2nQilPUCAwEAATANBgkqhkiG9w0BAQsF
AAOCAgEA2EM8rIIbsHCTDpm5Km5DUXk6kR0UcXWo2YacvHzEPilxTT1s8J6QrKNz
OkKGigEbP2I1WS73dxtc/58edmGjGawuCUN2GMRK/yvZUcUgsXJbtu+h4bHoi2ry
WfcbH9L9iHaS6XOPCBP4DZG4jawKAEi35jNra6x+JdFVPVmTQ08X7q7cZDyiIYEa
IplWjFodCwJA5QyL43hBMZYWfa1cRJQ+2uzEUt7w64sIhWCF7oSZUr50o4xS1MhB
ms9WdUqR9N+aySDv4SxUc1VO6E8YxBJhbPIIQJYnyRmD+gTz1JBLARLLrQ62Vjfw
juvQUcZQB0oXNPaEaDMS4clWP6s65zW2iZLsdSbRsTgsg1fvymuK5ZX94NkuHkMZ
zRPAyfwNVJuah33pKbUP4DUow4wwmg+AYu/iliWjxxZta+aI1GHZ6oLAV0+7lzD/
5t3MkW5OQhIw5OiRTHUwYDMpYoDoz4mIOGOTxWHhYR85zJAsh4gNIVD0Sk5mzPdp
Vx5ZTUDW1rS1JINGP009cihz3LvDSOiOx/Na0BS1MBZCmrMadBbDX6I99/qclCoF
AGiIEIopYhD4ASY8pBGvsAx3la5f/Sd9lANrDJbPlUMT1jlGXlo/2qaN3o/poqLc
nB0xUT1rwqXznSA7ZaevgGti6Zu14U3y0IIBTcjW1gVAXNr2w9s=
-----END CERTIFICATE-----
`

func TestCustomClient(t *testing.T) {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(caCrt))
	cliCrt, _ := tls.X509KeyPair([]byte(clientCrt), []byte(clientCsr))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	// a client with a cert for server check client
	client := &http.Client{Transport: tr}
	s, e := Get("https://httpbin.org/ip").WithCustomClient(client).String()
	if !strings.HasSuffix(e.Error(), ` x509: certificate signed by unknown authority`) {
		t.Fatalf(`want a error,but %v,%v`, s, e)
	}

	// use default client not check server cert
	s, e = Get("https://httpbin.org/ip").String()
	if e != nil || !strings.Contains(s, `"origin":`) {
		t.Fatalf(`want a "{"origin":"x.x.x.x"}",but %v,%v`, s, e)
	}
}
