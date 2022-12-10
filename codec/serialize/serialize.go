package serialize

type SerializeType string

const (
	SerializeTypeJson    SerializeType = "json"
	SerializeTypeMsgpack               = "msgpack"
)

var serializeMgr = map[SerializeType]Serialize{
	SerializeTypeJson:    &JsonSerialize{},
	SerializeTypeMsgpack: &MsgpackSerialize{},
}

type Serialize interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// RegisterSerialize 注册序列化方式
func RegisterSerialize(serializeType SerializeType, serialize Serialize) {
	serializeMgr[serializeType] = serialize
}

// GetSerialize 获取序列化方式
func GetSerialize(serializeType SerializeType) Serialize {
	return serializeMgr[serializeType]
}
