// Copyright 2020
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/astaxie/beego/pkg/infrastructure/config"
	"github.com/astaxie/beego/pkg/infrastructure/logs"
)

const etcdOpts = "etcdOpts"

type EtcdConfiger struct {
	prefix string
	client *clientv3.Client
	config.BaseConfiger
}

func newEtcdConfiger(client *clientv3.Client, prefix string) *EtcdConfiger {
	res := &EtcdConfiger{
		client: client,
		prefix: prefix,
	}

	res.BaseConfiger = config.NewBaseConfiger(res.reader)
	return res
}

// reader is an general implementation that read config from etcd.
func (e *EtcdConfiger) reader(ctx context.Context, key string) (string, error) {
	resp, err := get(e.client, ctx, e.prefix+key)
	if err != nil {
		return "", err
	}

	if resp.Count > 0 {
		return string(resp.Kvs[0].Value), nil
	}

	return "", nil
}

// Set do nothing and return an error
// I think write data to remote config center is not a good practice
func (e *EtcdConfiger) Set(ctx context.Context, key, val string) error {
	return errors.New("Unsupported operation")
}

// DIY return the original response from etcd
// be careful when you decide to use this
func (e *EtcdConfiger) DIY(ctx context.Context, key string) (interface{}, error) {
	return get(e.client, context.TODO(), key)
}

// GetSection in this implementation, we use section as prefix
func (e *EtcdConfiger) GetSection(ctx context.Context, section string) (map[string]string, error) {
	var (
		resp *clientv3.GetResponse
		err  error
	)

	if opts, ok := ctx.Value(etcdOpts).([]clientv3.OpOption); ok {
		opts = append(opts, clientv3.WithPrefix())
		resp, err = e.client.Get(context.TODO(), e.prefix+section, opts...)
	} else {
		resp, err = e.client.Get(context.TODO(), e.prefix+section, clientv3.WithPrefix())
	}

	if err != nil {
		return nil, errors.WithMessage(err, "GetSection failed")
	}
	res := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		res[string(kv.Key)] = string(kv.Value)
	}
	return res, nil
}

func (e *EtcdConfiger) SaveConfigFile(ctx context.Context, filename string) error {
	return errors.New("Unsupported operation")
}

// Unmarshaler is not very powerful because we lost the type information when we get configuration from etcd
// for example, when we got "5", we are not sure whether it's int 5, or it's string "5"
// TODO(support more complicated decoder)
func (e *EtcdConfiger) Unmarshaler(ctx context.Context, prefix string, obj interface{}, opt ...config.DecodeOption) error {
	res, err := e.GetSection(ctx, prefix)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("could not read config with prefix: %s", prefix))
	}

	prefixLen := len(e.prefix + prefix)
	m := make(map[string]string, len(res))
	for k, v := range res {
		m[k[prefixLen:]] = v
	}
	return mapstructure.Decode(m, obj)
}

// Sub return an sub configer.
func (e *EtcdConfiger) Sub(ctx context.Context, key string) (config.Configer, error) {
	return newEtcdConfiger(e.client, e.prefix+key), nil
}

// TODO remove this before release v2.0.0
func (e *EtcdConfiger) OnChange(ctx context.Context, key string, fn func(value string)) {

	buildOptsFunc := func() []clientv3.OpOption {
		if opts, ok := ctx.Value(etcdOpts).([]clientv3.OpOption); ok {
			opts = append(opts, clientv3.WithCreatedNotify())
			return opts
		}
		return []clientv3.OpOption{}
	}

	rch := e.client.Watch(ctx, e.prefix+key, buildOptsFunc()...)
	go func() {
		for {
			for resp := range rch {
				if err := resp.Err(); err != nil {
					logs.Error("listen to key but got error callback", err)
					break
				}

				for _, e := range resp.Events {
					if e.Kv == nil {
						continue
					}
					fn(string(e.Kv.Value))
				}
			}
			time.Sleep(time.Second)
			rch = e.client.Watch(ctx, e.prefix+key, buildOptsFunc()...)
		}
	}()

}

type EtcdConfigerProvider struct {
}

// Parse = ParseData([]byte(key))
// key must be json
func (provider *EtcdConfigerProvider) Parse(key string) (config.Configer, error) {
	return provider.ParseData([]byte(key))
}

// ParseData try to parse key as clientv3.Config, using this to build etcdClient
func (provider *EtcdConfigerProvider) ParseData(data []byte) (config.Configer, error) {
	cfg := &clientv3.Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "parse data to etcd config failed, please check your input")
	}

	cfg.DialOptions = []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	}
	client, err := clientv3.New(*cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "create etcd client failed")
	}

	return newEtcdConfiger(client, ""), nil
}

func get(client *clientv3.Client, ctx context.Context, key string) (*clientv3.GetResponse, error) {
	var (
		resp *clientv3.GetResponse
		err  error
	)
	if opts, ok := ctx.Value(etcdOpts).([]clientv3.OpOption); ok {
		resp, err = client.Get(ctx, key, opts...)
	} else {
		resp, err = client.Get(ctx, key)
	}

	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("read config from etcd with key %s failed", key))
	}
	return resp, err
}

func WithEtcdOption(ctx context.Context, opts ...clientv3.OpOption) context.Context {
	return context.WithValue(ctx, etcdOpts, opts)
}

func init() {
	config.Register("json", &EtcdConfigerProvider{})
}
