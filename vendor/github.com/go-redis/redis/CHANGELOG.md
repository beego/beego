# Changelog

## v6.13

- Ring got new options called `HashReplicas` and `Hash`. It is recommended to set `HashReplicas = 1000` for better keys distribution between shards.
- Cluster client was optimized to use much less memory when reloading cluster state.
- PubSub.ReceiveMessage is re-worked to not use ReceiveTimeout so it does not lose data when timeout occurres. In most cases it is recommended to use PubSub.Channel instead.
- Dialer.KeepAlive is set to 5 minutes by default.

## v6.12

- ClusterClient got new option called `ClusterSlots` which allows to build cluster of normal Redis Servers that don't have cluster mode enabled. See https://godoc.org/github.com/go-redis/redis#example-NewClusterClient--ManualSetup
