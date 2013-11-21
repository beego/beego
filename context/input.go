package context

import (
	"bytes"
	"github.com/astaxie/beego/session"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type BeegoInput struct {
	CruSession  session.SessionStore
	Param       map[string]string
	req         *http.Request
	RequestBody []byte
}

func NewInput(req *http.Request) *BeegoInput {
	return &BeegoInput{
		Param: make(map[string]string),
		req:   req,
	}
}

func (input *BeegoInput) Protocol() string {
	return input.req.Proto
}

func (input *BeegoInput) Uri() string {
	return input.req.RequestURI
}

func (input *BeegoInput) Url() string {
	return input.req.URL.String()
}

func (input *BeegoInput) Site() string {
	return input.Scheme() + "://" + input.Domain()
}

func (input *BeegoInput) Scheme() string {
	if input.req.URL.Scheme != "" {
		return input.req.URL.Scheme
	} else if input.req.TLS == nil {
		return "http"
	} else {
		return "https"
	}
}

func (input *BeegoInput) Domain() string {
	return input.Host()
}

func (input *BeegoInput) Host() string {
	if input.req.Host != "" {
		hostParts := strings.Split(input.req.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return input.req.Host
	}
	return "localhost"
}

func (input *BeegoInput) Method() string {
	return input.req.Method
}

func (input *BeegoInput) Is(method string) bool {
	return input.Method() == method
}

func (input *BeegoInput) IsAjax() bool {
	return input.Header("X-Requested-With") == "XMLHttpRequest"
}

func (input *BeegoInput) IsSecure() bool {
	return input.Scheme() == "https"
}

func (input *BeegoInput) IsWebsocket() bool {
	return input.Header("Upgrade") == "websocket"
}

func (input *BeegoInput) IsUpload() bool {
	return input.req.MultipartForm != nil
}

func (input *BeegoInput) IP() string {
	ips := input.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	ip := strings.Split(input.req.RemoteAddr, ":")
	if len(ip) > 0 {
		return ip[0]
	}
	return "127.0.0.1"
}

func (input *BeegoInput) Proxy() []string {
	if ips := input.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (input *BeegoInput) Refer() string {
	return input.Header("Referer")
}

func (input *BeegoInput) SubDomains() string {
	parts := strings.Split(input.Host(), ".")
	return strings.Join(parts[len(parts)-2:], ".")
}

func (input *BeegoInput) Port() int {
	parts := strings.Split(input.req.Host, ":")
	if len(parts) == 2 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}

func (input *BeegoInput) UserAgent() string {
	return input.Header("User-Agent")
}

func (input *BeegoInput) Params(key string) string {
	if v, ok := input.Param[key]; ok {
		return v
	}
	return ""
}

func (input *BeegoInput) Query(key string) string {
	return input.req.Form.Get(key)
}

func (input *BeegoInput) Header(key string) string {
	return input.req.Header.Get(key)
}

func (input *BeegoInput) Cookie(key string) string {
	ck, err := input.req.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

func (input *BeegoInput) Session(key interface{}) interface{} {
	return input.CruSession.Get(key)
}

func (input *BeegoInput) Body() []byte {
	requestbody, _ := ioutil.ReadAll(input.req.Body)
	input.req.Body.Close()
	bf := bytes.NewBuffer(requestbody)
	input.req.Body = ioutil.NopCloser(bf)
	input.RequestBody = requestbody
	return requestbody
}
