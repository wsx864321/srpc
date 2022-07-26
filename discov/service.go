package discov

type Service struct {
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints"`
}

type Endpoint struct {
	//InstanceID string `json:"instance_id"`
	ServiceName string `json:"service_name"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Weight      int    `json:"weight"`
	//Protocol   protocol.ProtocolType `json:"protocol"`
	Enable bool `json:"enable"`
}
