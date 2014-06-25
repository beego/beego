// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package beego

import (
	"fmt"
	"net/url"
	"strings"
)

// FlashData is a tools to maintain data when using across request.
type FlashData struct {
	Data map[string]string
}

// NewFlash return a new empty FlashData struct.
func NewFlash() *FlashData {
	return &FlashData{
		Data: make(map[string]string),
	}
}

// Notice writes notice message to flash.
func (fd *FlashData) Notice(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["notice"] = msg
	} else {
		fd.Data["notice"] = fmt.Sprintf(msg, args...)
	}
}

// Warning writes warning message to flash.
func (fd *FlashData) Warning(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["warning"] = msg
	} else {
		fd.Data["warning"] = fmt.Sprintf(msg, args...)
	}
}

// Error writes error message to flash.
func (fd *FlashData) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		fd.Data["error"] = msg
	} else {
		fd.Data["error"] = fmt.Sprintf(msg, args...)
	}
}

// Store does the saving operation of flash data.
// the data are encoded and saved in cookie.
func (fd *FlashData) Store(c *Controller) {
	c.Data["flash"] = fd.Data
	var flashValue string
	for key, value := range fd.Data {
		flashValue += "\x00" + key + "\x23" + FlashSeperator + "\x23" + value + "\x00"
	}
	c.Ctx.SetCookie(FlashName, url.QueryEscape(flashValue), 0, "/")
}

// ReadFromRequest parsed flash data from encoded values in cookie.
func ReadFromRequest(c *Controller) *FlashData {
	flash := NewFlash()
	if cookie, err := c.Ctx.Request.Cookie(FlashName); err == nil {
		v, _ := url.QueryUnescape(cookie.Value)
		vals := strings.Split(v, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, "\x23"+FlashSeperator+"\x23")
				if len(kv) == 2 {
					flash.Data[kv[0]] = kv[1]
				}
			}
		}
		//read one time then delete it
		c.Ctx.SetCookie(FlashName, "", -1, "/")
	}
	c.Data["flash"] = flash.Data
	return flash
}
