package serialize

type SerializeType string

const (
	SerializeTypeJson    SerializeType = "json"
	SerializeTypeMsgpack               = "msgpack"
	SerializeTypeProto                 = "proto"
)

var serializeMgr = map[SerializeType]Serialize{
	SerializeTypeJson:    NewJsonSerialize(),
	SerializeTypeMsgpack: NewMsgpackSerialize(),
	SerializeTypeProto:   NewProtoSerialize(),
}

type Serialize interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// RegisterSerialize 注册序列化方式,可以自定义序列化方式
func RegisterSerialize(serializeType SerializeType, serialize Serialize) {
	serializeMgr[serializeType] = serialize
}

// GetSerialize 获取序列化方式
func GetSerialize(serializeType SerializeType) Serialize {
	return serializeMgr[serializeType]
}
