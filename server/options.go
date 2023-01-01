package server

import (
	"time"

	"github.com/wsx864321/srpc/logger"

	"github.com/wsx864321/srpc/discov/etcd"

	"github.com/wsx864321/srpc/discov"

	"github.com/wsx864321/srpc/codec/serialize"
	"github.com/wsx864321/srpc/transport"
)

var defaultOptions = &Options{
	ip:           "0.0.0.0",
	port:         9557,
	network:      transport.NetworkTCP,
	serialize:    serialize.SerializeTypeJson,
	timeout:      5 * time.Second,
	writeTimeout: 1 * time.Second,
	discovery:    etcd.NewETCDRegister(),
	logger:       logger.NewSweetLog(),
}

type Options struct {
	ip           string
	port         int
	network      transport.Transport
	serialize    serialize.SerializeType
	timeout      time.Duration
	writeTimeout time.Duration
	discovery    discov.Discovery
	logger       logger.Log
}

type Option func(opt *Options)

func WithIP(ip string) Option {
	return func(opt *Options) {
		opt.ip = ip
	}
}

func WithPort(port int) Option {
	return func(opt *Options) {
		opt.port = port
	}
}

func WithNetWork(network transport.Transport) Option {
	return func(opt *Options) {
		opt.network = network
	}
}

func WithSerialize(serializeType serialize.SerializeType) Option {
	return func(opt *Options) {
		opt.serialize = serializeType
	}
}

func WithTimeOut(duration time.Duration) Option {
	return func(opt *Options) {
		opt.timeout = duration
	}
}

func WithWriteTimeout(duration time.Duration) Option {
	return func(opt *Options) {
		opt.writeTimeout = duration
	}
}

func WithDiscovery(discovery discov.Discovery) Option {
	return func(opt *Options) {
		opt.discovery = discovery
	}
}

func WithLogger(log logger.Log) Option {
	return func(opt *Options) {
		opt.logger = log
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
