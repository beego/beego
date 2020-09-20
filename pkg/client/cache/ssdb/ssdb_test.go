package ssdb

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/astaxie/beego/pkg/client/cache"
)

func TestSsdbcacheCache(t *testing.T) {

	ssdbAddr := os.Getenv("SSDB_ADDR")
	if ssdbAddr == "" {
		ssdbAddr = "127.0.0.1:8888"
	}

	cc, err := cache.NewCache("ssdb", fmt.Sprintf(`{"conn": "%s"}`, ssdbAddr))
	if err != nil {
		t.Error("init err")
	}

	// test put and exist
	if _, err := cc.IsExist("test_key"); err != nil {
		t.Error("check err")
	}
	timeoutDuration := 10 * time.Second
	//timeoutDuration := -10*time.Second   if timeoutDuration is negtive,it means permanent
	if err = cc.Put("test_key", "test_val", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	// test put and exist
	b, err := cc.IsExist("test_key")
	if err != nil {
		t.Error("check err")
	}
	if b == false {
		t.Error("check err")
	}

	// Get test done
	if err = cc.Put("test_key", "test_val", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if v, _ := cc.Get("test_key"); v != "test_val" {
		t.Error("get Error")
	}

	//inc/dec test done
	if err = cc.Put("test_key", "2", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if _, err = cc.Incr("test_key"); err != nil {
		t.Error("incr Error", err)
	}

	v, _ := cc.Get("test_key")
	if v, err := strconv.Atoi(v.(string)); err != nil || v != 3 {
		t.Error("get err")
	}

	if _, err = cc.Decr("test_key"); err != nil {
		t.Error("decr error")
	}

	// test del
	if err = cc.Put("test_key", "3", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	v, _ = cc.Get("test_key")
	if v, err := strconv.Atoi(v.(string)); err != nil || v != 3 {
		t.Error("get err")
	}
	if err := cc.Delete("test_key"); err == nil {
		if ok, _ := cc.IsExist("test_key"); ok {
			t.Error("delete err")
		}
	}

	//test string
	if err = cc.Put("test_key", "test_val", -10*time.Second); err != nil {
		t.Error("set Error", err)
	}
	if ok, _ := cc.IsExist("test_key"); ok {
		t.Error("check err")
	}
	if v, _ := cc.Get("test_key"); v.(string) != "test_val" {
		t.Error("get err")
	}

	//test GetMulti done
	if err = cc.Put("k1", "v1", -10*time.Second); err != nil {
		t.Error("set Error", err)
	}
	if ok, _ := cc.IsExist("k1"); !ok {
		t.Error("check err")
	}
	vv, err := cc.GetMulti([]string{"k1", "k2"})
	if len(vv) != 2 {
		t.Error("getmulti error")
	}
	if vv[0].(string) != "v1" {
		t.Error("getmulti error")
	}
	if vv[1].(string) != "v1" {
		t.Error("getmulti error")
	}

	// test clear all done
	if err = cc.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}
