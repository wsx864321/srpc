package server

import (
	"time"

	"github.com/wsx864321/sweet_rpc/codec/serialize"
	"github.com/wsx864321/sweet_rpc/transport"
)

var defaultOptions = &Options{
	Addr:         "0.0.0.0:9577",
	Protocol:     transport.ProtocolTCP,
	Serialize:    serialize.SerializeTypeJson,
	ReadTimeout:  20 * time.Second,
	WriteTimeout: 20 * time.Second,
}

type Options struct {
	Addr         string
	Protocol     transport.Protocol
	Serialize    serialize.SerializeType
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Option func(opt *Options)

func WithAddr(addr string) Option {
	return func(opt *Options) {
		opt.Addr = addr
	}
}

func WithProtocol(protocol transport.Protocol) Option {
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

// NewOptions 初始化option
func NewOptions(opts ...Option) *Options {
	opt := defaultOptions
	for _, fn := range opts {
		fn(opt)
	}

	return opt
}
