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
	IP:           "0.0.0.0",
	Port:         9557,
	Network:      transport.NetworkTCP,
	Serialize:    serialize.SerializeTypeJson,
	Timeout:      5 * time.Second,
	WriteTimeout: 1 * time.Second,
	Discovery:    etcd.NewETCDRegister(),
	Logger:       logger.NewSweetLog(),
}

type Options struct {
	IP           string
	Port         int
	Network      transport.Transport
	Serialize    serialize.SerializeType
	Timeout      time.Duration
	WriteTimeout time.Duration
	Discovery    discov.Discovery
	Logger       logger.Log
}

type Option func(opt *Options)

func WithIP(ip string) Option {
	return func(opt *Options) {
		opt.IP = ip
	}
}

func WithPort(port int) Option {
	return func(opt *Options) {
		opt.Port = port
	}
}

func WithNetWork(network transport.Transport) Option {
	return func(opt *Options) {
		opt.Network = network
	}
}

func WithSerialize(serializeType serialize.SerializeType) Option {
	return func(opt *Options) {
		opt.Serialize = serializeType
	}
}

func WithTimeOut(duration time.Duration) Option {
	return func(opt *Options) {
		opt.Timeout = duration
	}
}

func WithWriteTimeout(duration time.Duration) Option {
	return func(opt *Options) {
		opt.WriteTimeout = duration
	}
}

func WithDiscovery(discovery discov.Discovery) Option {
	return func(opt *Options) {
		opt.Discovery = discovery
	}
}

func WithLogger(log logger.Log) Option {
	return func(opt *Options) {
		opt.Logger = log
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
