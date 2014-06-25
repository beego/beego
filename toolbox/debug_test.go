// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package toolbox

import (
	"testing"
)

type mytype struct {
	next *mytype
	prev *mytype
}

func TestPrint(t *testing.T) {
	Display("v1", 1, "v2", 2, "v3", 3)
}

func TestPrintPoint(t *testing.T) {
	var v1 = new(mytype)
	var v2 = new(mytype)

	v1.prev = nil
	v1.next = v2

	v2.prev = v1
	v2.next = nil

	Display("v1", v1, "v2", v2)
}

func TestPrintString(t *testing.T) {
	str := GetDisplayString("v1", 1, "v2", 2)
	println(str)
}
