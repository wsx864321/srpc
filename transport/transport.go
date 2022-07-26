package transport

type Transport string

const (
	ProtocolTCP  Transport = "tcp"
	ProtocolUDP            = "udp"
	ProtocolKCP            = "kcp"
	ProtocolHTTP           = "http"
	ProtocolQUIC           = "quic"
)
