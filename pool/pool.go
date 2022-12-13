package pool

import (
	"net"
	"sync"
)

type ConnectionPool struct {
}

type pool struct {
	conns *sync.Map
}

type connection struct {
	net.Conn
}
