package redis

import (
	"crypto/tls"
	"time"
)

// UniversalOptions information is required by UniversalClient to establish
// connections.
type UniversalOptions struct {
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Addrs []string

	// The sentinel master name.
	// Only failover clients.
	MasterName string

	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int

	// Only cluster clients.

	// Enables read only queries on slave nodes.
	ReadOnly bool

	MaxRedirects   int
	RouteByLatency bool

	// Common options

	OnConnect          func(*Conn) error
	MaxRetries         int
	Password           string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	TLSConfig          *tls.Config
}

func (o *UniversalOptions) cluster() *ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:6379"}
	}

	return &ClusterOptions{
		Addrs:          o.Addrs,
		MaxRedirects:   o.MaxRedirects,
		RouteByLatency: o.RouteByLatency,
		ReadOnly:       o.ReadOnly,

		OnConnect:          o.OnConnect,
		MaxRetries:         o.MaxRetries,
		Password:           o.Password,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
}

func (o *UniversalOptions) failover() *FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}

	return &FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		DB:            o.DB,

		OnConnect:          o.OnConnect,
		MaxRetries:         o.MaxRetries,
		Password:           o.Password,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
}

func (o *UniversalOptions) simple() *Options {
	addr := "127.0.0.1:6379"
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &Options{
		Addr: addr,
		DB:   o.DB,

		OnConnect:          o.OnConnect,
		MaxRetries:         o.MaxRetries,
		Password:           o.Password,
		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
}

// --------------------------------------------------------------------

// UniversalClient is an abstract client which - based on the provided options -
// can connect to either clusters, or sentinel-backed failover instances or simple
// single-instance servers. This can be useful for testing cluster-specific
// applications locally.
type UniversalClient interface {
	Cmdable
	Watch(fn func(*Tx) error, keys ...string) error
	Process(cmd Cmder) error
	WrapProcess(fn func(oldProcess func(cmd Cmder) error) func(cmd Cmder) error)
	Subscribe(channels ...string) *PubSub
	PSubscribe(channels ...string) *PubSub
	Close() error
}

var _ UniversalClient = (*Client)(nil)
var _ UniversalClient = (*ClusterClient)(nil)

// NewUniversalClient returns a new multi client. The type of client returned depends
// on the following three conditions:
//
// 1. if a MasterName is passed a sentinel-backed FailoverClient will be returned
// 2. if the number of Addrs is two or more, a ClusterClient will be returned
// 3. otherwise, a single-node redis Client will be returned.
func NewUniversalClient(opts *UniversalOptions) UniversalClient {
	if opts.MasterName != "" {
		return NewFailoverClient(opts.failover())
	} else if len(opts.Addrs) > 1 {
		return NewClusterClient(opts.cluster())
	}
	return NewClient(opts.simple())
}
