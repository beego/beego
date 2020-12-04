package cache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/client/orm"
)

type Data struct {
	Id  int
	Key string `orm:"unique"`
	Val int
}

var DBARGS = struct {
	Driver string
	Source string
	Debug  string
}{
	os.Getenv("ORM_DRIVER"),
	os.Getenv("ORM_SOURCE"),
	os.Getenv("ORM_DEBUG"),
}

func init() {
	orm.RegisterModel(new(Data))
	if err := orm.RegisterDataBase("default", DBARGS.Driver, DBARGS.Source); err != nil {
		panic(fmt.Sprintf("can not register database: %v", err))
	}
	if err := orm.RunSyncdb("default", false, false); err != nil {
		panic(fmt.Sprintf("can not run sync db: %v", err))
	}
}

type dbReader struct {
	ormer orm.Ormer
}

func NewDbReader(orm orm.Ormer) *dbReader {
	return &dbReader{ormer: orm}
}

func (r *dbReader) Query(key string) (interface{}, error) {
	kv := &Data{Key: key}
	err := r.ormer.Read(kv, "key")
	if err != nil {
		return nil, err
	}
	return kv.Val, nil
}

type dbWriter struct {
	ormer orm.Ormer
}

func NewDbWriter(orm orm.Ormer) *dbWriter {
	return &dbWriter{ormer: orm}
}

func (w *dbWriter) Update(key string, val interface{}) error {
	valInt, ok := val.(int)
	if !ok {
		return fmt.Errorf("val can't convert to int")
	}
	kv := &Data{Key: key, Val: valInt}
	if _, err := w.ormer.InsertOrUpdate(kv); err != nil {
		return err
	}
	return nil
}

func TestCacheAside(t *testing.T) {
	o := orm.NewOrm()
	c, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Errorf("init cache failed, err: %s", err.Error())
		return
	}
	r := NewDbReader(o)
	w := NewDbWriter(o)
	ca := NewCacheAside(c, r, w, 10*time.Second)
	ctx := context.Background()

	key1 := "a"
	val1 := 1

	actualVal, _ := ca.Get(ctx, key1)
	if actualVal != nil {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}

	if _, err := o.Insert(&Data{Key: key1, Val: val1}); err != nil {
		t.Errorf("insert failed, err: %s", err.Error())
		return
	}
	defer func() {
		if _, err := o.Delete(&Data{Key: key1}, "key"); err != nil {
			t.Errorf("delte failed, err: %s", err.Error())
			return
		}
	}()

	// check if the data is get from data source successfully
	actualVal, _ = ca.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}

	// check if the data is put to cache successfully
	actualVal, _ = c.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}

	// check if the data in cache and data source is updated successfully after execute CacheAside.Put
	val1 = 0
	if err := ca.Put(ctx, key1, val1, 10*time.Second); err != nil {
		t.Errorf("put data failed, err: %s", err.Error())
		return
	}
	actualVal, err = r.Query(key1)
	if err != nil {
		t.Errorf("query data from datasource failed, err: %s", err.Error())
		return
	}
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}
	// the update of data will delete the cache, so the first query from cache will cause the cache miss,
	// then the updated val will put to cache
	if ok, _ := c.IsExist(ctx, key1); ok {
		t.Errorf("key isn't delte from cache successfully")
		return
	}
	actualVal, _ = ca.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}
	actualVal, _ = c.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}

	// check if the data in cache and data source is updated successfully after execute CacheAside.Incr
	val1++
	if err := ca.Incr(ctx, key1); err != nil {
		t.Errorf("incr data failed, err: %s", err.Error())
		return
	}
	actualVal, err = r.Query(key1)
	if err != nil {
		t.Errorf("query data from datasource failed, err: %s", err.Error())
		return
	}
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}
	// the update of data will delete the cache, so the first query from cache will cause the cache miss,
	// then the updated val will put to cache
	if ok, _ := c.IsExist(ctx, key1); ok {
		t.Errorf("key isn't delte from cache successfully")
		return
	}
	actualVal, _ = ca.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}
	actualVal, _ = c.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}

	// check if the data in cache and data source is updated successfully after execute CacheAside.Decr
	val1--
	if err := ca.Decr(ctx, key1); err != nil {
		t.Errorf("incr data failed, err: %s", err.Error())
		return
	}
	actualVal, err = r.Query(key1)
	if err != nil {
		t.Errorf("query data from datasource failed, err: %s", err.Error())
		return
	}
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}
	// the update of data will delete the cache, so the first query from cache will cause the cache miss,
	// then the updated val will put to cache
	if ok, _ := c.IsExist(ctx, key1); ok {
		t.Errorf("key isn't delte from cache successfully")
		return
	}
	actualVal, _ = ca.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}
	actualVal, _ = c.Get(ctx, key1)
	if actualVal == nil || actualVal.(int) != val1 {
		t.Errorf("check actualVal failed, expect %d, but %v actually", val1, actualVal)
		return
	}

	// test GetMulti
	key2 := "b"
	val2 := 2
	if err := ca.Put(ctx, key2, val2, 10*time.Second); err != nil {
		t.Errorf("put failed, err: %s", err.Error())
		return
	}
	defer func() {
		if _, err := o.Delete(&Data{Key: key2}, "key"); err != nil {
			t.Errorf("delte failed, err: %s", err.Error())
			return
		}
	}()

	vals, err := ca.GetMulti(ctx, []string{key1, key2})
	if err != nil {
		t.Errorf("getmulti failed, err: %s", err.Error())
		return
	}
	v, ok := vals[0].(int)
	if !ok || v != val1 {
		t.Errorf("getmulti failed, expect %d, but %v actually", val1, v)
		return
	}
	v, ok = vals[1].(int)
	if !ok || v != val2 {
		t.Errorf("getmulti failed, expect %d, but %v actually", val2, v)
		return
	}

	vals, err = ca.GetMulti(ctx, []string{key1, key2, "notexist"})
	if err == nil {
		t.Error("getmulti failed")
		return
	}
	v, ok = vals[0].(int)
	if !ok || v != val1 {
		t.Errorf("getmulti failed, expect %d, but %v actually", val1, v)
		return
	}
	v, ok = vals[1].(int)
	if !ok || v != val2 {
		t.Errorf("getmulti failed, expect %d, but %v actually", val2, v)
		return
	}
	if vals[2] != nil {
		t.Error("getmulti failed")
		return
	}

	// test incr&decr when the key isn't exist
	if err := ca.Incr(ctx, "notexist"); err == nil {
		t.Error("incr failed")
		return
	}
	if err := ca.Decr(ctx, "notexist"); err == nil {
		t.Error("decr failed")
		return
	}
}
