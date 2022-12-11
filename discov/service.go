package discov

type Service struct {
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints"`
}

// Endpoint 增加序列化和传输协议字段（序列化协议交由server端去定义是更加合理的）
type Endpoint struct {
	//InstanceID string `json:"instance_id"`
	ServiceName string `json:"service_name"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Weight      int    `json:"weight"`
	Network     string `json:"network"`
	Serialize   string `json:"serialize"`
	Enable      bool   `json:"enable"`
}
