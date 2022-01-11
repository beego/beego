package session

type ProviderType string

const (
	ProviderCookie        ProviderType = `cookie`
	ProviderFile          ProviderType = `file`
	ProviderMemory        ProviderType = `memory`
	ProviderCouchbase     ProviderType = `couchbase`
	ProviderLedis         ProviderType = `ledis`
	ProviderMemcache      ProviderType = `memcache`
	ProviderMysql         ProviderType = `mysql`
	ProviderPostgresql    ProviderType = `postgresql`
	ProviderRedis         ProviderType = `redis`
	ProviderRedisCluster  ProviderType = `redis_cluster`
	ProviderRedisSentinel ProviderType = `redis_sentinel`
	ProviderSsdb          ProviderType = `ssdb`
)
