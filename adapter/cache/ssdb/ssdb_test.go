package ssdb

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/beego/beego/v2/adapter/cache"
)

const (
	initError = "init err"
	setError = "set Error"
	checkError = "check err"
	getError = "get err"
	getMultiError = "GetMulti Error"
)

func TestSsdbcacheCache(t *testing.T) {
	ssdbAddr := os.Getenv("SSDB_ADDR")
	if ssdbAddr == "" {
		ssdbAddr = "127.0.0.1:8888"
	}

	ssdb, err := cache.NewCache("ssdb", fmt.Sprintf(`{"conn": "%s"}`, ssdbAddr))
	if err != nil {
		t.Error(initError)
	}

	// test put and exist
	if ssdb.IsExist("ssdb") {
		t.Error(checkError)
	}
	timeoutDuration := 10 * time.Second
	// timeoutDuration := -10*time.Second   if timeoutDuration is negtive,it means permanent
	if err = ssdb.Put("ssdb", "ssdb", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if !ssdb.IsExist("ssdb") {
		t.Error(checkError)
	}

	// Get test done
	if err = ssdb.Put("ssdb", "ssdb", timeoutDuration); err != nil {
		t.Error(setError, err)
	}

	if v := ssdb.Get("ssdb"); v != "ssdb" {
		t.Error("get Error")
	}

	// inc/dec test done
	if err = ssdb.Put("ssdb", "2", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if err = ssdb.Incr("ssdb"); err != nil {
		t.Error("incr Error", err)
	}

	if v, err := strconv.Atoi(ssdb.Get("ssdb").(string)); err != nil || v != 3 {
		t.Error(getError)
	}

	if err = ssdb.Decr("ssdb"); err != nil {
		t.Error("decr error")
	}

	// test del
	if err = ssdb.Put("ssdb", "3", timeoutDuration); err != nil {
		t.Error(setError, err)
	}
	if v, err := strconv.Atoi(ssdb.Get("ssdb").(string)); err != nil || v != 3 {
		t.Error(getError)
	}
	if err := ssdb.Delete("ssdb"); err == nil {
		if ssdb.IsExist("ssdb") {
			t.Error("delete err")
		}
	}

	// test string
	if err = ssdb.Put("ssdb", "ssdb", -10*time.Second); err != nil {
		t.Error(setError, err)
	}
	if !ssdb.IsExist("ssdb") {
		t.Error(checkError)
	}
	if v := ssdb.Get("ssdb").(string); v != "ssdb" {
		t.Error(getError)
	}

	// test GetMulti done
	if err = ssdb.Put("ssdb1", "ssdb1", -10*time.Second); err != nil {
		t.Error(setError, err)
	}
	if !ssdb.IsExist("ssdb1") {
		t.Error(checkError)
	}
	vv := ssdb.GetMulti([]string{"ssdb", "ssdb1"})
	if len(vv) != 2 {
		t.Error(getMultiError)
	}
	if vv[0].(string) != "ssdb" {
		t.Error(getMultiError)
	}
	if vv[1].(string) != "ssdb1" {
		t.Error(getMultiError)
	}

	// test clear all done
	if err = ssdb.ClearAll(); err != nil {
		t.Error("clear all err")
	}
	if ssdb.IsExist("ssdb") || ssdb.IsExist("ssdb1") {
		t.Error(checkError)
	}
}
