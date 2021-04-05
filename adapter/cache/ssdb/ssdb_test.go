package ssdb

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	assert.Nil(t, err)

	assert.False(t, ssdb.IsExist("ssdb"))
	// test put and exist
	timeoutDuration := 3 * time.Second
	// timeoutDuration := -10*time.Second   if timeoutDuration is negtive,it means permanent
	assert.Nil(t, ssdb.Put("ssdb", "ssdb", timeoutDuration))
	assert.True(t, ssdb.IsExist("ssdb"))

	assert.Nil(t, ssdb.Put("ssdb", "ssdb", timeoutDuration))

	assert.Equal(t, "ssdb", ssdb.Get("ssdb"))

	// inc/dec test done
	assert.Nil(t, ssdb.Put("ssdb", "2", timeoutDuration))

	assert.Nil(t, ssdb.Incr("ssdb"))

	v, err := strconv.Atoi(ssdb.Get("ssdb").(string))
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	assert.Nil(t, ssdb.Decr("ssdb"))

	assert.Nil(t, ssdb.Put("ssdb", "3", timeoutDuration))

	// test del
	v, err = strconv.Atoi(ssdb.Get("ssdb").(string))
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	assert.Nil(t, ssdb.Delete("ssdb"))
	assert.False(t, ssdb.IsExist("ssdb"))

	// test string
	assert.Nil(t, ssdb.Put("ssdb", "ssdb", -10*time.Second))

	assert.True(t, ssdb.IsExist("ssdb"))
	assert.Equal(t, "ssdb", ssdb.Get("ssdb"))

	// test GetMulti done
	assert.Nil(t, ssdb.Put("ssdb1", "ssdb1", -10*time.Second))
	assert.True(t, ssdb.IsExist("ssdb1") )

	vv := ssdb.GetMulti([]string{"ssdb", "ssdb1"})
	assert.Equal(t, 2, len(vv))

	assert.Equal(t, "ssdb", vv[0])
	assert.Equal(t, "ssdb1", vv[1])

	assert.Nil(t, ssdb.ClearAll())
	assert.False(t, ssdb.IsExist("ssdb"))
	assert.False(t, ssdb.IsExist("ssdb1"))
	// test clear all done
}
