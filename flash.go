package beego

import (
	"fmt"
	"net/url"
	"strings"
)

// the separation string when encoding flash data.
const BEEGO_FLASH_SEP = "#BEEGOFLASH#"

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
		flashValue += "\x00" + key + BEEGO_FLASH_SEP + value + "\x00"
	}
	c.Ctx.SetCookie("BEEGO_FLASH", url.QueryEscape(flashValue), 0, "/")
}

// ReadFromRequest parsed flash data from encoded values in cookie.
func ReadFromRequest(c *Controller) *FlashData {
	flash := &FlashData{
		Data: make(map[string]string),
	}
	if cookie, err := c.Ctx.Request.Cookie("BEEGO_FLASH"); err == nil {
		v, _ := url.QueryUnescape(cookie.Value)
		vals := strings.Split(v, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, BEEGO_FLASH_SEP)
				if len(kv) == 2 {
					flash.Data[kv[0]] = kv[1]
				}
			}
		}
		//read one time then delete it
		c.Ctx.SetCookie("BEEGO_FLASH", "", -1, "/")
	}
	c.Data["flash"] = flash.Data
	return flash
}
