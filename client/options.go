package client

import (
	"github.com/wsx864321/sweet_rpc/interceptor"
	"github.com/wsx864321/sweet_rpc/lb"
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

// NewOptions 初始化option
func NewOptions(opts ...Option) *Options {
	opt := defaultOptions
	for _, fn := range opts {
		fn(opt)
	}

	return opt
}
