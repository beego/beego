// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package logs

import (
	"testing"
)

func TestConn(t *testing.T) {
	log := NewLogger(1000)
	log.SetLogger("conn", `{"net":"tcp","addr":":7020"}`)
	log.Info("info")
}
