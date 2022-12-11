package serialize

import "encoding/json"

type JsonSerialize struct {
}

func NewJsonSerialize() Serialize {
	return &JsonSerialize{}
}

func (j *JsonSerialize) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonSerialize) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
