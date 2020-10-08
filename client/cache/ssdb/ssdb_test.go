package ssdb

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/astaxie/beego/client/cache"
)

func TestSsdbcacheCache(t *testing.T) {

	ssdbAddr := os.Getenv("SSDB_ADDR")
	if ssdbAddr == "" {
		ssdbAddr = "127.0.0.1:8888"
	}

	ssdb, err := cache.NewCache("ssdb", fmt.Sprintf(`{"conn": "%s"}`, ssdbAddr))
	if err != nil {
		t.Error("init err")
	}

	// test put and exist
	if res, _ := ssdb.IsExist(context.Background(), "ssdb"); res {
		t.Error("check err")
	}
	timeoutDuration := 10 * time.Second
	// timeoutDuration := -10*time.Second   if timeoutDuration is negtive,it means permanent
	if err = ssdb.Put(context.Background(), "ssdb", "ssdb", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := ssdb.IsExist(context.Background(), "ssdb"); !res {
		t.Error("check err")
	}

	// Get test done
	if err = ssdb.Put(context.Background(), "ssdb", "ssdb", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if v, _ := ssdb.Get(context.Background(), "ssdb"); v != "ssdb" {
		t.Error("get Error")
	}

	// inc/dec test done
	if err = ssdb.Put(context.Background(), "ssdb", "2", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if err = ssdb.Incr(context.Background(), "ssdb"); err != nil {
		t.Error("incr Error", err)
	}

	val, _ := ssdb.Get(context.Background(), "ssdb")
	if v, err := strconv.Atoi(val.(string)); err != nil || v != 3 {
		t.Error("get err")
	}

	if err = ssdb.Decr(context.Background(), "ssdb"); err != nil {
		t.Error("decr error")
	}

	// test del
	if err = ssdb.Put(context.Background(), "ssdb", "3", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	val, _ = ssdb.Get(context.Background(), "ssdb")
	if v, err := strconv.Atoi(val.(string)); err != nil || v != 3 {
		t.Error("get err")
	}
	if err := ssdb.Delete(context.Background(), "ssdb"); err == nil {
		if e, _ := ssdb.IsExist(context.Background(), "ssdb"); e {
			t.Error("delete err")
		}
	}

	// test string
	if err = ssdb.Put(context.Background(), "ssdb", "ssdb", -10*time.Second); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := ssdb.IsExist(context.Background(), "ssdb"); !res {
		t.Error("check err")
	}
	if v, _ := ssdb.Get(context.Background(), "ssdb"); v.(string) != "ssdb" {
		t.Error("get err")
	}

	// test GetMulti done
	if err = ssdb.Put(context.Background(), "ssdb1", "ssdb1", -10*time.Second); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := ssdb.IsExist(context.Background(), "ssdb1"); !res {
		t.Error("check err")
	}
	vv, _ := ssdb.GetMulti(context.Background(), []string{"ssdb", "ssdb1"})
	if len(vv) != 2 {
		t.Error("getmulti error")
	}
	if vv[0].(string) != "ssdb" {
		t.Error("getmulti error")
	}
	if vv[1].(string) != "ssdb1" {
		t.Error("getmulti error")
	}

	// test clear all done
	if err = ssdb.ClearAll(context.Background()); err != nil {
		t.Error("clear all err")
	}
	e1, _ := ssdb.IsExist(context.Background(), "ssdb")
	e2, _ := ssdb.IsExist(context.Background(), "ssdb1")
	if e1 || e2 {
		t.Error("check err")
	}
}
