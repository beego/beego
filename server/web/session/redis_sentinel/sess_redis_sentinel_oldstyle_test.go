package redis_sentinel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProvider_initOldStyle(t *testing.T) {
	cp := &Provider{}
	cp.initOldStyle("127.0.0.1:26379;127.0.0.2:26379,100,1qaz2wsx,1,mymaster,30,60,3")
	assert.Equal(t, "127.0.0.1:26379;127.0.0.2:26379", cp.SavePath)
	assert.Equal(t, 100, cp.Poolsize)
	assert.Equal(t, "1qaz2wsx", cp.Password)
	assert.Equal(t, 1, cp.DbNum)
	assert.Equal(t, "mymaster", cp.MasterName)
	assert.Equal(t, 30*time.Second, cp.idleTimeout)
	assert.Equal(t, 60*time.Second, cp.idleCheckFrequency)
	assert.Equal(t, 3, cp.MaxRetries)
}
