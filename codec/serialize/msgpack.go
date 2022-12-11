package serialize

import "github.com/vmihailenco/msgpack"

type MsgpackSerialize struct {
}

func NewMsgpackSerialize() Serialize {
	return &MsgpackSerialize{}
}

func (m *MsgpackSerialize) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m *MsgpackSerialize) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
