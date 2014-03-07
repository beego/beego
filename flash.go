package beego

import (
	"fmt"
	"net/url"
)

// FlashData is a tools to maintain data when using across request.
type FlashData struct {
	Data      map[string]string
	Name      string
	Seperator string
}

// NewFlash return a new empty FlashData struct.
func NewFlash() *FlashData {
	return &FlashData{
		Data:      make(map[string]string),
		Name:      FlashName,
		Seperator: FlashSeperator,
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
		flashValue += "\x00" + key + "\x23" + fd.Seperator + "\x23" + value + "\x00"
	}
	c.Ctx.SetCookie(fd.Name, url.QueryEscape(flashValue), 0, "/")
}
