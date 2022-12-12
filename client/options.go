package client

import (
	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/lb"
	"time"
)

var defaultOptions = &Options{
	writeTimeout: 5 * time.Second,
	readTimeout:  5 * time.Second,
	lb:           lb.NewRandom(),
}

type Options struct {
	serviceName  string
	interceptors []interceptor.ClientInterceptor
	writeTimeout time.Duration
	timeout      time.Duration
	readTimeout  time.Duration
	lb           lb.LoadBalance
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

func WithLoadBalance(balance lb.LoadBalance) Option {
	return func(o *Options) {
		o.lb = balance
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
