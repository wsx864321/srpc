package transport

import "net"

type Transport string

const (
	ProtocolTCP  Transport = "tcp"
	ProtocolUDP            = "udp"
	ProtocolKCP            = "kcp"
	ProtocolHTTP           = "http"
	ProtocolQUIC           = "quic"
)

type ServerTransport interface {
	Listen(addr string) (net.Listener, error)
}

type ServerTransportMgr struct {
	transports map[Transport]ServerTransport
}

func NewServerTransportMgr() *ServerTransportMgr {
	return &ServerTransportMgr{
		map[Transport]ServerTransport{
			ProtocolTCP: newTCP(),
		},
	}
}

func (s *ServerTransportMgr) Gen(transport Transport) ServerTransport {
	if val, ok := s.transports[transport]; ok {
		return val
	}
	return s.transports[ProtocolTCP]
}
