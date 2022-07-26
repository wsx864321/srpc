package codec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodec_Encode(t *testing.T) {
	c := NewCodec()
	resp, err := c.Encode(1, 1, 11111, []byte("testService"), []byte("index"), []byte("matedata"), []byte("payload"))
	assert.NoError(t, err)

	fmt.Println(resp)
}

func TestCodec_decodeHeader(t *testing.T) {
	h := Header{
		17, 1, 1, 1, 11, 5, 8, 7, 11111,
	}
	c := NewCodec()
	req, _ := c.Encode(1, 1, 11111, []byte("testService"), []byte("index"), []byte("matedata"), []byte("payload"))

	resp, err := c.decodeHeader(req[:fixedBytesConst])

	assert.NoError(t, err)

	assert.Equal(t, h, *resp)
}

func TestCodec_Decode(t *testing.T) {
	c := NewCodec()
	req, _ := c.Encode(1, 1, 11111, []byte("testService"), []byte("index"), []byte("matedata"), []byte("payload"))

	resp, err := c.Decode(req)

	assert.NoError(t, err)

	assert.Equal(t, resp.ServiceName, "testService")
	assert.Equal(t, resp.ServiceMethod, "index")
	assert.Equal(t, resp.MetaData, []byte("matedata"))
	assert.Equal(t, resp.Payload, []byte("payload"))
}
