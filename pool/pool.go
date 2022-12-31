package pool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed     = errors.New("pool is closed")
	ErrGetTimeout = errors.New("poo get client timeout")
)

// Pool 基本方法
//type Pool interface {
//	Get(address, network string) (net.Conn, error)
//
//	Put(conn net.Conn) error
//
//	Close(conn net.Conn) error
//
//	Release()
//
//	Len() int
//}

type Pool struct {
	opts  *Options
	conns *sync.Map
	mu    *sync.Mutex
}

// NewPool todo 增加异步client链接检查，防止下游节点因为服务下线而连接还保存在map中造成缓慢的内存泄露和fd的浪费
func NewPool(opts ...Option) *Pool {
	opt := NewOptions(opts...)
	if !(opt.initialCap <= opt.maxCap) {
		panic("invalid capacity settings")
	}
	if opt.factory == nil {
		panic("invalid factory func settings")
	}
	if opt.close == nil {
		panic("invalid close func settings")
	}

	return &Pool{
		opts:  opt,
		conns: &sync.Map{},
		mu:    &sync.Mutex{},
	}
}

// Get ...
func (p *Pool) Get(ctx context.Context, network, address string) (net.Conn, error) {
	if value, ok := p.conns.Load(p.getKey(network, address)); ok {
		if cp, ok := value.(*channelPool); ok {
			conn, err := cp.Get(ctx, network, address)
			return conn, err
		}
	}

	// check - lock -check 单例
	p.mu.Lock()
	defer p.mu.Unlock()
	if value, ok := p.conns.Load(p.getKey(network, address)); ok {
		if cp, ok := value.(*channelPool); ok {
			conn, err := cp.Get(ctx, network, address)
			return conn, err
		}
	}

	cp, err := NewChannelPool(
		WithInitialCap(p.opts.initialCap),
		WithMaxCap(p.opts.maxCap),
		WithFactory(p.opts.factory),
		WithClose(p.opts.close),
		WithPing(p.opts.ping),
		WithIdleTimeout(p.opts.idleTimeout),
		WithNetwork(network),
		WithAddress(address),
		WithDailTimeout(p.opts.dailTimeout),
	)
	if err != nil {
		return nil, err
	}

	p.conns.Store(p.getKey(network, address), cp)

	return cp.Get(ctx, network, address)
}

func (p *Pool) Put(network, address string, conn net.Conn) {
	if value, ok := p.conns.Load(p.getKey(network, address)); ok {
		if cp, ok := value.(*channelPool); ok {
			cp.Put(conn)
		}
	}
}

func (p *Pool) CloseAll() {
	p.conns.Range(func(key, value interface{}) bool {
		if cp, ok := value.(*channelPool); ok {
			cp.Release()
		}

		return true
	})
}

func (p *Pool) getKey(network, address string) string {
	return fmt.Sprintf("%s://%s", network, address)
}
