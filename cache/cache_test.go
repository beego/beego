package cache

import (
	"testing"
	"time"
)

func Test_cache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	if err = bm.Put("astaxie", 1, 10); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("astaxie") {
		t.Error("check err")
	}

	if err = bm.Put("astaxie", 1, 10); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("astaxie"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("astaxie"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("astaxie")
	if bm.IsExist("astaxie") {
		t.Error("delete err")
	}
}
