package redis

import (
	"crypto/tls"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/internal"
	"github.com/go-redis/redis/internal/pool"
)

//------------------------------------------------------------------------------

// FailoverOptions are used to configure a failover client and should
// be passed to NewFailoverClient.
type FailoverOptions struct {
	// The master name.
	MasterName string
	// A seed list of host:port addresses of sentinel nodes.
	SentinelAddrs []string

	// Following options are copied from Options struct.

	OnConnect func(*Conn) error

	Password string
	DB       int

	MaxRetries int

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration

	TLSConfig *tls.Config
}

func (opt *FailoverOptions) options() *Options {
	return &Options{
		Addr: "FailoverClient",

		OnConnect: opt.OnConnect,

		DB:       opt.DB,
		Password: opt.Password,

		MaxRetries: opt.MaxRetries,

		DialTimeout:  opt.DialTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,

		PoolSize:           opt.PoolSize,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,

		TLSConfig: opt.TLSConfig,
	}
}

// NewFailoverClient returns a Redis client that uses Redis Sentinel
// for automatic failover. It's safe for concurrent use by multiple
// goroutines.
func NewFailoverClient(failoverOpt *FailoverOptions) *Client {
	opt := failoverOpt.options()
	opt.init()

	failover := &sentinelFailover{
		masterName:    failoverOpt.MasterName,
		sentinelAddrs: failoverOpt.SentinelAddrs,

		opt: opt,
	}

	c := Client{
		baseClient: baseClient{
			opt:      opt,
			connPool: failover.Pool(),

			onClose: func() error {
				return failover.Close()
			},
		},
	}
	c.baseClient.init()
	c.setProcessor(c.Process)

	return &c
}

//------------------------------------------------------------------------------

type SentinelClient struct {
	baseClient
}

func NewSentinelClient(opt *Options) *SentinelClient {
	opt.init()
	c := &SentinelClient{
		baseClient: baseClient{
			opt:      opt,
			connPool: newConnPool(opt),
		},
	}
	c.baseClient.init()
	return c
}

func (c *SentinelClient) PubSub() *PubSub {
	pubsub := &PubSub{
		opt: c.opt,

		newConn: func(channels []string) (*pool.Conn, error) {
			return c.newConn()
		},
		closeConn: c.connPool.CloseConn,
	}
	pubsub.init()
	return pubsub
}

func (c *SentinelClient) GetMasterAddrByName(name string) *StringSliceCmd {
	cmd := NewStringSliceCmd("SENTINEL", "get-master-addr-by-name", name)
	c.Process(cmd)
	return cmd
}

func (c *SentinelClient) Sentinels(name string) *SliceCmd {
	cmd := NewSliceCmd("SENTINEL", "sentinels", name)
	c.Process(cmd)
	return cmd
}

type sentinelFailover struct {
	sentinelAddrs []string

	opt *Options

	pool     *pool.ConnPool
	poolOnce sync.Once

	mu          sync.RWMutex
	masterName  string
	_masterAddr string
	sentinel    *SentinelClient
}

func (d *sentinelFailover) Close() error {
	return d.resetSentinel()
}

func (d *sentinelFailover) Pool() *pool.ConnPool {
	d.poolOnce.Do(func() {
		d.opt.Dialer = d.dial
		d.pool = newConnPool(d.opt)
	})
	return d.pool
}

func (d *sentinelFailover) dial() (net.Conn, error) {
	addr, err := d.MasterAddr()
	if err != nil {
		return nil, err
	}
	return net.DialTimeout("tcp", addr, d.opt.DialTimeout)
}

func (d *sentinelFailover) MasterAddr() (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	addr, err := d.masterAddr()
	if err != nil {
		return "", err
	}
	d._switchMaster(addr)

	return addr, nil
}

func (d *sentinelFailover) masterAddr() (string, error) {
	// Try last working sentinel.
	if d.sentinel != nil {
		addr, err := d.sentinel.GetMasterAddrByName(d.masterName).Result()
		if err == nil {
			addr := net.JoinHostPort(addr[0], addr[1])
			return addr, nil
		}

		internal.Logf("sentinel: GetMasterAddrByName name=%q failed: %s",
			d.masterName, err)
		d._resetSentinel()
	}

	for i, sentinelAddr := range d.sentinelAddrs {
		sentinel := NewSentinelClient(&Options{
			Addr: sentinelAddr,

			DialTimeout:  d.opt.DialTimeout,
			ReadTimeout:  d.opt.ReadTimeout,
			WriteTimeout: d.opt.WriteTimeout,

			PoolSize:    d.opt.PoolSize,
			PoolTimeout: d.opt.PoolTimeout,
			IdleTimeout: d.opt.IdleTimeout,
		})

		masterAddr, err := sentinel.GetMasterAddrByName(d.masterName).Result()
		if err != nil {
			internal.Logf("sentinel: GetMasterAddrByName master=%q failed: %s",
				d.masterName, err)
			sentinel.Close()
			continue
		}

		// Push working sentinel to the top.
		d.sentinelAddrs[0], d.sentinelAddrs[i] = d.sentinelAddrs[i], d.sentinelAddrs[0]
		d.setSentinel(sentinel)

		addr := net.JoinHostPort(masterAddr[0], masterAddr[1])
		return addr, nil
	}

	return "", errors.New("redis: all sentinels are unreachable")
}

func (c *sentinelFailover) switchMaster(addr string) {
	c.mu.Lock()
	c._switchMaster(addr)
	c.mu.Unlock()
}

func (c *sentinelFailover) _switchMaster(addr string) {
	if c._masterAddr == addr {
		return
	}

	internal.Logf("sentinel: new master=%q addr=%q",
		c.masterName, addr)
	_ = c.Pool().Filter(func(cn *pool.Conn) bool {
		return cn.RemoteAddr().String() != addr
	})
	c._masterAddr = addr
}

func (d *sentinelFailover) setSentinel(sentinel *SentinelClient) {
	d.discoverSentinels(sentinel)
	d.sentinel = sentinel
	go d.listen(sentinel)
}

func (d *sentinelFailover) resetSentinel() error {
	var err error
	d.mu.Lock()
	if d.sentinel != nil {
		err = d._resetSentinel()
	}
	d.mu.Unlock()
	return err
}

func (d *sentinelFailover) _resetSentinel() error {
	err := d.sentinel.Close()
	d.sentinel = nil
	return err
}

func (d *sentinelFailover) discoverSentinels(sentinel *SentinelClient) {
	sentinels, err := sentinel.Sentinels(d.masterName).Result()
	if err != nil {
		internal.Logf("sentinel: Sentinels master=%q failed: %s", d.masterName, err)
		return
	}
	for _, sentinel := range sentinels {
		vals := sentinel.([]interface{})
		for i := 0; i < len(vals); i += 2 {
			key := vals[i].(string)
			if key == "name" {
				sentinelAddr := vals[i+1].(string)
				if !contains(d.sentinelAddrs, sentinelAddr) {
					internal.Logf(
						"sentinel: discovered new sentinel=%q for master=%q",
						sentinelAddr, d.masterName,
					)
					d.sentinelAddrs = append(d.sentinelAddrs, sentinelAddr)
				}
			}
		}
	}
}

func (d *sentinelFailover) listen(sentinel *SentinelClient) {
	pubsub := sentinel.PubSub()
	defer pubsub.Close()

	err := pubsub.Subscribe("+switch-master")
	if err != nil {
		internal.Logf("sentinel: Subscribe failed: %s", err)
		d.resetSentinel()
		return
	}

	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			if err == pool.ErrClosed {
				d.resetSentinel()
				return
			}
			internal.Logf("sentinel: ReceiveMessage failed: %s", err)
			continue
		}

		switch msg.Channel {
		case "+switch-master":
			parts := strings.Split(msg.Payload, " ")
			if parts[0] != d.masterName {
				internal.Logf("sentinel: ignore addr for master=%q", parts[0])
				continue
			}
			addr := net.JoinHostPort(parts[3], parts[4])
			d.switchMaster(addr)
		}
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
