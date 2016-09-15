package utils

import (
	"net/http"
	"strings"
)

// 获取 客户端 IP
func GetRequestIP(r *http.Request) string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		if arr := strings.Split(ips, ","); len(arr) > 0 && arr[0] != "" {
			rip := strings.Split(arr[0], ":")
			return rip[0]
		}
	}

	if ip := strings.Split(r.RemoteAddr, ":"); len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}
