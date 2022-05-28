package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	magicNum   uint8 = 0x11
	version    uint8 = 0x1
	fixedBytes       = 24
)

type Header struct {
	MagicNum          uint8
	Version           uint8 // version
	MsgType           uint8 // msg type e.g. :   0x0: general req,  0x1: heartbeat
	CompressType      uint8 // compression or not :  0x0: not compression,  0x1: compression
	ServiceNameSize   uint16
	ServiceMethodSize uint16
	MateSize          uint32
	DataSize          uint32
	Seq               uint64 // stream ID
}

/**
协议设计
|MagicNum|Version|MsgType|CompressType|ServiceNameSize|ServiceMethodSize|MateSize|DataSize|Seq     |ServiceName|ServiceMethod|MetaData|Payload
|1byte   |1byte  |1byte  |1byte       |2byte          |2byte            |4byte   |4byte   |8byte   |
*/

type Message struct {
	*Header
	ServiceName   string
	ServiceMethod string
	MetaData      []byte
	Payload       []byte
}

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
}

// Encode ...
func (c *Codec) Encode(msgType, compressType, uint8, seq uint64, serviceName, serviceMethod, metaData, payload []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, fixedBytes+len(serviceName)+len(serviceMethod)+len(metaData)+len(payload)))
	if err := binary.Write(buffer, binary.BigEndian, magicNum); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, magicNum); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, version); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, msgType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, compressType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, len(serviceName)); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, len(serviceMethod)); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, len(metaData)); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, len(payload)); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, seq); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, serviceName); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, serviceMethod); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, metaData); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *Codec) Decode(data []byte) (*Message, error) {
	if len(data) < fixedBytes {
		return nil, errors.New("data length is at least 23")
	}

	return nil, nil
}

func (c *Codec) decodeHeader(data []byte) (*Header, error) {
	var (
		magicNum          uint8
		version           uint8
		msgType           uint8
		compressType      uint8
		serviceNameSize   uint16
		serviceMethodSize uint16
		mateSize          uint32
		dataSize          uint32
		seq               uint64
	)
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &magicNum); err != nil {
		return nil, err
	}

	if magicNum != fixedBytes {
		return nil, errors.New("invalid data")
	}

	if err := binary.Read(buffer, binary.BigEndian, &version); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &msgType); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &compressType); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &serviceNameSize); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &serviceMethodSize); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &mateSize); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &dataSize); err != nil {
		return nil, err
	}

	if err := binary.Read(buffer, binary.BigEndian, &seq); err != nil {
		return nil, err
	}

	return &Header{
		magicNum,
		version,
		msgType,
		compressType,
		serviceNameSize,
		serviceMethodSize,
		mateSize,
		dataSize,
		seq,
	}, nil
}
