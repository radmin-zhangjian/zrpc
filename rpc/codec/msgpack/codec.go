package msgpack

import (
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	"zrpc/rpc/zio"
)

type codec struct {
	conn net.Conn
}

func New(conn net.Conn) *codec {
	return &codec{conn: conn}
}

// Encoder 编码
func (c *codec) Encoder(data zio.Response) ([]byte, error) {
	// msgpack 编码
	buf, err := msgpack.Marshal(data)
	if err != nil {
		log.Fatalln("Failed to encode:", err)
		return buf, err
	}
	return buf, nil
}

// Decoder 解码
func (c *codec) Decoder(b []byte) (zio.Response, error) {
	// msgpack 解码
	var data zio.Response
	err := msgpack.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Failed to decode:", err)
		return data, err
	}
	return data, nil
}
