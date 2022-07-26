package server

import (
	"time"

	"github.com/wsx864321/sweet_rpc/discov/etcd"

	"github.com/wsx864321/sweet_rpc/discov"

	"github.com/wsx864321/sweet_rpc/codec/serialize"
	"github.com/wsx864321/sweet_rpc/transport"
)

var defaultOptions = &Options{
	IP:           "0.0.0.0",
	Port:         9557,
	Protocol:     transport.ProtocolTCP,
	Serialize:    serialize.SerializeTypeJson,
	ReadTimeout:  20 * time.Second,
	WriteTimeout: 20 * time.Second,
	Discovery:    etcd.NewETCDRegister(),
}

type Options struct {
	IP           string
	Port         int
	Protocol     transport.Transport
	Serialize    serialize.SerializeType
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Discovery    discov.Discovery
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

func WithProtocol(protocol transport.Transport) Option {
	return func(opt *Options) {
		opt.Protocol = protocol
	}
}

func WithSerialize(serializeType serialize.SerializeType) Option {
	return func(opt *Options) {
		opt.Serialize = serializeType
	}
}

func WithReadTimeout(duration time.Duration) Option {
	return func(opt *Options) {
		opt.ReadTimeout = duration
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

// NewOptions 初始化option
func NewOptions(opts ...Option) *Options {
	opt := defaultOptions
	for _, fn := range opts {
		fn(opt)
	}

	return opt
}
