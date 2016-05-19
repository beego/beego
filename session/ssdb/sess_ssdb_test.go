package ssdb

import (
	"fmt"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	p := &SsdbProvider{}
	p.SessionInit(300, "127.0.0.1:8888")
	if p.host != "127.0.0.1" || p.port != 8888 {
		t.Error("host:port err")
	}
	if p.client == nil {
		t.Error("client err")
	}
	ss, err := p.SessionRead("1")
	if err != nil {
		t.Error(err)
	}
	err = ss.Set("key", "value")
	if err != nil {
		t.Error(err)
	}
	if ss.Get("key") != "value" {
		t.Error("Get err")
	}
	err = ss.Delete("key")
	//err = ss.Flush()
	if err != nil {
		t.Error(err)
	}
	if ss.Get("key") == "value" {
		t.Error("Delete/Flush err")
	}
	if ss.SessionID() != "1" {
		t.Error("id err")
	}

	ss.Set("key1", "value1")
	var w http.ResponseWriter
	ss.SessionRelease(w)
	new, e := p.SessionRead("1")
	if new == nil || e != nil {
		t.Error(e)
	}
	if !p.SessionExist("1") {
		t.Error("SessionExist err")
	}
	newS, er := p.SessionRegenerate("1", "3")
	if er != nil || newS == nil {
		t.Error("SessionRegenerate err")
	}
	if p.SessionExist("1") {
		t.Error("SessionExist err")
	}
	fmt.Println(newS)

}
