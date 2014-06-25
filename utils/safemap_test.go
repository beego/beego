// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package utils

import (
	"testing"
)

func Test_beemap(t *testing.T) {
	bm := NewBeeMap()
	if !bm.Set("astaxie", 1) {
		t.Error("set Error")
	}
	if !bm.Check("astaxie") {
		t.Error("check err")
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	bm.Delete("astaxie")
	if bm.Check("astaxie") {
		t.Error("delete err")
	}
}
