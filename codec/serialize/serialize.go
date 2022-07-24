package serialize

type SerializeType string

const (
	SerializeTypeJson    SerializeType = "json"
	SerializeTypeMsgpack               = "msgpack"
)

type Serialize interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}
