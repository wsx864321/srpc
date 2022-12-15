package pool

import (
	"errors"
	"net"
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
)

// Pool 基本方法
type Pool interface {
	Get(address, network string) (net.Conn, error)

	Put(conn net.Conn) error

	Close(conn net.Conn) error

	Release()

	Len() int
}
