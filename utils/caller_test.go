// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package utils

import (
	"strings"
	"testing"
)

func TestGetFuncName(t *testing.T) {
	name := GetFuncName(TestGetFuncName)
	t.Log(name)
	if !strings.HasSuffix(name, ".TestGetFuncName") {
		t.Error("get func name error")
	}
}
