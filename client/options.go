package client

import (
	"github.com/wsx864321/srpc/discov"
	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/lb"
	"github.com/wsx864321/srpc/pool"
	"time"
)

var defaultOptions = &Options{
	writeTimeout: 500 * time.Millisecond,
	timeout:      5 * time.Second,
	readTimeout:  500 * time.Millisecond,
	lbName:       lb.LoadBalanceRandom,
	dis:          etcd.NewETCDRegister(),
	pool:         pool.NewPool(),
}

type Options struct {
	serviceName  string
	interceptors []interceptor.ClientInterceptor
	writeTimeout time.Duration
	timeout      time.Duration
	readTimeout  time.Duration
	lbName       string
	dis          discov.Discovery
	pool         *pool.Pool
}

type Option func(o *Options)

func WithServiceName(srvName string) Option {
	return func(o *Options) {
		o.serviceName = srvName
	}
}

func WithInterceptors(interceptor ...interceptor.ClientInterceptor) Option {
	return func(o *Options) {
		o.interceptors = interceptor
	}
}

func WithWriteTimeout(duration time.Duration) Option {
	return func(o *Options) {
		o.writeTimeout = duration
	}
}

func WithTimeout(duration time.Duration) Option {
	return func(o *Options) {
		o.timeout = duration
	}
}

func WithReadTimeout(duration time.Duration) Option {
	return func(o *Options) {
		o.readTimeout = duration
	}
}

func WithLoadBalance(name string) Option {
	return func(o *Options) {
		o.lbName = name
	}
}

func WithDiscovery(discovery discov.Discovery) Option {
	return func(o *Options) {
		o.dis = discovery
	}
}

func WithPool(pool *pool.Pool) Option {
	return func(o *Options) {
		o.pool = pool
	}
}

// NewOptions 初始化option
func NewOptions(opts ...Option) *Options {
	opt := defaultOptions
	for _, fn := range opts {
		fn(opt)
	}

	return opt
}
