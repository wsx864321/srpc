package transport

type Protocol string

const (
	ProtocolTCP  Protocol = "tcp"
	ProtocolUDP           = "udp"
	ProtocolKCP           = "kcp"
	ProtocolHTTP          = "http"
	ProtocolQUIC          = "quic"
)
