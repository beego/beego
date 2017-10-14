package session

import (
	"testing"
	"encoding/json"
)

func TestRedisConfig(t *testing.T) {
	op := RedisOptions{
		Addr:     "192.168.0.23:6379",
		Password: "123456",
		PoolSize: 100,
	}
	b, _ := json.Marshal(op)
	t.Log(string(b))

}
