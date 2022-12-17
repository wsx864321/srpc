package pool

import (
	"net"
	"time"
)

var defaultOptions = &Options{
	initialCap: 1,
	maxCap:     100,
	factory: func(network, address string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout(network, address, timeout)
	},
	close: func(conn net.Conn) error {
		return conn.Close()
	},
	ping:        nil,
	idleTimeout: 1 * time.Minute,
	network:     "tcp",
	address:     "127.0.0.1:7777",
	dailTimeout: 100 * time.Millisecond,
}

type Options struct {
	//连接池中拥有的最小连接数
	initialCap int
	//最大并发存活连接数
	maxCap int
	//生成连接的方法
	factory func(network, address string, timeout time.Duration) (net.Conn, error)
	//关闭连接的方法
	close func(conn net.Conn) error
	//检查连接是否有效的方法
	ping func(conn net.Conn) error
	//连接最大空闲时间，超过该事件则将失效
	idleTimeout time.Duration
	// network eg:tcp udp
	network string
	// ip+port 0.0.0.0:9999
	address string
	// dailTimeout
	dailTimeout time.Duration
}

type Option func(opts *Options)

func WithInitialCap(cap int) Option {
	return func(opts *Options) {
		opts.initialCap = cap
	}
}

func WithMaxCap(cap int) Option {
	return func(opts *Options) {
		opts.maxCap = cap
	}
}

func WithFactory(factory func(network, address string, timeout time.Duration) (net.Conn, error)) Option {
	return func(opts *Options) {
		opts.factory = factory
	}
}

func WithClose(close func(conn net.Conn) error) Option {
	return func(opts *Options) {
		opts.close = close
	}
}

func WithPing(ping func(conn net.Conn) error) Option {
	return func(opts *Options) {
		opts.ping = ping
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(opts *Options) {
		opts.idleTimeout = idleTimeout
	}
}

func WithNetwork(network string) Option {
	return func(opts *Options) {
		opts.network = network
	}
}

func WithAddress(addr string) Option {
	return func(opts *Options) {
		opts.address = addr
	}
}

func WithDailTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.dailTimeout = timeout
	}
}

func NewOptions(opts ...Option) *Options {
	opt := defaultOptions
	for _, o := range opts {
		o(opt)
	}

	return opt
}
