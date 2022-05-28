package discov

type Service struct {
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints"`
}

type Endpoint struct {
	//InstanceID string                `json:"instance_id"`
	ServerName string `json:"server_name"`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	//Protocol   protocol.ProtocolType `json:"protocol"`
	Enable bool `json:"enable"`
}
