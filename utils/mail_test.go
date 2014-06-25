// Beego (http://beego.me/)

// @description beego is an open-source, high-performance web framework for the Go programming language.

// @link        http://github.com/astaxie/beego for the canonical source repository

// @license     http://github.com/astaxie/beego/blob/master/LICENSE

// @authors     astaxie

package utils

import "testing"

func TestMail(t *testing.T) {
	config := `{"username":"astaxie@gmail.com","password":"astaxie","host":"smtp.gmail.com","port":587}`
	mail := NewEMail(config)
	if mail.Username != "astaxie@gmail.com" {
		t.Fatal("email parse get username error")
	}
	if mail.Password != "astaxie" {
		t.Fatal("email parse get password error")
	}
	if mail.Host != "smtp.gmail.com" {
		t.Fatal("email parse get host error")
	}
	if mail.Port != 587 {
		t.Fatal("email parse get port error")
	}
	mail.To = []string{"xiemengjun@gmail.com"}
	mail.From = "astaxie@gmail.com"
	mail.Subject = "hi, just from beego!"
	mail.Text = "Text Body is, of course, supported!"
	mail.HTML = "<h1>Fancy Html is supported, too!</h1>"
	mail.AttachFile("/Users/astaxie/github/beego/beego.go")
	mail.Send()
}
