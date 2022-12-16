package transport

import "net"

type udp struct {
}

func newUDP() ServerTransport {
	return &udp{}
}

func (u *udp) Listen(addr string) (net.Listener, error) {
	return net.Listen("udp", addr)
}
