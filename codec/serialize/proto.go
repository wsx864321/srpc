package serialize

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type ProtoSerialize struct {
}

func NewProtoSerialize() Serialize {
	return &ProtoSerialize{}
}

func (p *ProtoSerialize) Marshal(v interface{}) ([]byte, error) {
	message, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("v is not implement proto.Message,v type:%v", reflect.TypeOf(v))
	}

	return proto.Marshal(message)
}

func (p *ProtoSerialize) Unmarshal(data []byte, v interface{}) error {
	message, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("v is not implement proto.Message,v type:%v", reflect.TypeOf(v))
	}

	return proto.Unmarshal(data, message)
}
