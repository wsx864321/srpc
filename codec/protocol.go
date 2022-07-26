package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	magicNumConst   uint8 = 0x11
	versionConst    uint8 = 0x1
	fixedBytesConst       = 24
)

/**
协议设计,TLV,|tag|length|value
|MagicNum|Version|MsgType|CompressType|ServiceNameSize|ServiceMethodSize|MateSize|DataSize|Seq     |ServiceName|ServiceMethod|MetaData|Payload
|1byte   |1byte  |1byte  |1byte       |2byte          |2byte            |4byte   |4byte   |8byte   |x bytes    |x bytes      |x bytes |x bytes
*/

// todo 修改成rpc调用四元素 caller callerFunc callee callFunc
type Message struct {
	*Header
	ServiceName   string
	ServiceMethod string
	MetaData      []byte
	Payload       []byte
}

// Header 设计上考虑了内存对齐
type Header struct {
	MagicNum          uint8
	Version           uint8 // version
	MsgType           uint8 // msg type e.g. :   0x0: general req,  0x1: heartbeat
	CompressType      uint8 // compression or not :  0x0: not compression,  0x1: compression
	ServiceNameSize   uint16
	ServiceMethodSize uint16
	MetaSize          uint32
	PayloadSize       uint32
	Seq               uint64 // stream ID
}

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
}

// Encode ...
func (c *Codec) Encode(msgType, compressType uint8, seq uint64, serviceName, serviceMethod, metaData, payload []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, fixedBytesConst+len(serviceName)+len(serviceMethod)+len(metaData)+len(payload)))
	if err := binary.Write(buffer, binary.BigEndian, magicNumConst); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, versionConst); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, msgType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, compressType); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint16(len(serviceName))); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint16(len(serviceMethod))); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint32(len(metaData))); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, uint32(len(payload))); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.BigEndian, seq); err != nil {
		return nil, err
	}

	if _, err := buffer.Write(serviceName); err != nil {
		return nil, err
	}
	if _, err := buffer.Write(serviceMethod); err != nil {
		return nil, err
	}

	if _, err := buffer.Write(metaData); err != nil {
		return nil, err
	}

	if _, err := buffer.Write(payload); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *Codec) Decode(data []byte) (*Message, error) {
	if len(data) < fixedBytesConst {
		return nil, errors.New("data length is at least 23")
	}

	var msg Message
	header, err := c.decodeHeader(data[:fixedBytesConst])
	if err != nil {
		return nil, err
	}

	serviceNameStartPos := fixedBytesConst
	msg.ServiceName = string(data[serviceNameStartPos : serviceNameStartPos+int(header.ServiceNameSize)])

	serviceMethodStartPos := serviceNameStartPos + int(header.ServiceNameSize)
	msg.ServiceMethod = string(data[serviceMethodStartPos : serviceMethodStartPos+int(header.ServiceMethodSize)])

	metaDataStartPos := serviceMethodStartPos + int(header.ServiceMethodSize)
	msg.MetaData = data[metaDataStartPos : metaDataStartPos+int(header.MetaSize)]

	payloadStartPos := metaDataStartPos + int(header.MetaSize)
	msg.Payload = data[payloadStartPos : payloadStartPos+int(header.PayloadSize)]

	msg.Header = header

	return &msg, nil
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

	if magicNum != magicNum {
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
