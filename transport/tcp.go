package transport

import "net"

type tcp struct {
}

func newTCP() ServerTransport {
	return &tcp{}
}

func (t *tcp) Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
