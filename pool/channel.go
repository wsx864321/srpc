package pool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	//"reflect"
)

var (
	//ErrMaxActiveConnReached 连接池超限
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
)

// channelPool 存放连接信息
type channelPool struct {
	*Options

	mu           sync.RWMutex
	conns        chan *idleConn
	openingConns int
}

type idleConn struct {
	conn net.Conn
	t    time.Time
}

// NewChannelPool 初始化连接
func NewChannelPool(opts ...Option) (*channelPool, error) {
	opt := NewOptions(opts...)
	c := &channelPool{
		Options:      opt,
		mu:           sync.RWMutex{},
		conns:        make(chan *idleConn, opt.maxCap),
		openingConns: opt.initialCap,
	}

	for i := 0; i < opt.initialCap; i++ {
		if c.factory == nil {
			fmt.Println(0)
			fmt.Println(opts)
			fmt.Println(NewOptions(opts...))
		}
		conn, err := c.factory(opt.network, opt.address, opt.dailTimeout)
		if err != nil {
			c.release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}

	return c, nil
}

// getConns 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	return c.conns
}

// Get 从pool中取一个连接 todo 增加超时控制
func (c *channelPool) get(ctx context.Context, network, address string) (net.Conn, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ErrGetTimeout
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该连接
					c.close(wrapConn.conn)
					continue
				}
			}
			//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
			if c.pingFunc != nil {
				if err := c.ping(wrapConn.conn); err != nil {
					c.close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:

			if c.openingConns >= c.maxCap {
				for {
					select {
					case <-ctx.Done():
						return nil, ErrGetTimeout
					case wrapConn := <-conns:
						if wrapConn == nil {
							return nil, ErrClosed
						}
						//判断是否超时，超时则丢弃
						if timeout := c.idleTimeout; timeout > 0 {
							if wrapConn.t.Add(timeout).Before(time.Now()) {
								//丢弃并关闭该连接
								c.close(wrapConn.conn)
								continue
							}
						}
						//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
						if c.pingFunc != nil {
							if err := c.ping(wrapConn.conn); err != nil {
								c.close(wrapConn.conn)
								continue
							}
						}
						return wrapConn.conn, nil
					}
				}
			}

			conn, err := c.factory(network, address, c.dailTimeout)
			if err != nil {
				return nil, err
			}

			c.mu.Lock()
			defer c.mu.Unlock()

			c.openingConns++
			return conn, nil
		}
	}
}

// Put 将连接放回pool中
func (c *channelPool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	if c.conns == nil {
		return c.close(conn)
	}

	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该连接
		return c.close(conn)
	}

}

// Close 关闭单条连接
func (c *channelPool) close(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	if c.closeFunc == nil {
		return nil
	}
	c.mu.Lock()
	c.openingConns--
	c.mu.Unlock()

	return c.close(conn)
}

// Ping 检查单条连接是否有效
func (c *channelPool) ping(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.pingFunc(conn)
}

// release 释放连接池中所有连接
func (c *channelPool) release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.pingFunc = nil
	closeFun := c.close
	c.closeFunc = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		//log.Printf("Type %v\n",reflect.TypeOf(wrapConn.conn))
		closeFun(wrapConn.conn)
	}
}

// Len 连接池中已有的连接
func (c *channelPool) len() int {
	return len(c.getConns())
}
