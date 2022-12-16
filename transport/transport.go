package transport

import "net"

type Transport string

const (
	NetworkTCP  Transport = "tcp"
	NetworkUDP            = "udp"
	NetworkKCP            = "kcp"
	NetworkHTTP           = "http"
	NetworkQUIC           = "quic"
)

var transportMgr = map[Transport]ServerTransport{
	NetworkTCP: newTCP(),
	NetworkUDP: newUDP(),
}

type ServerTransport interface {
	Listen(addr string) (net.Listener, error)
}

// GetTransport 获取传输协议
func GetTransport(transport Transport) ServerTransport {
	return transportMgr[transport]
}

// RegisterTransport 注册传输协议
func RegisterTransport(transport Transport, serverTransport ServerTransport) {
	transportMgr[transport] = serverTransport
}
