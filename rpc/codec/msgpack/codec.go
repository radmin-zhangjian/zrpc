package msgpack

import (
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	cdc "zrpc/rpc/codec"
)

type codec struct {
	conn net.Conn
}

func New(conn net.Conn) *codec {
	return &codec{conn: conn}
}

func FuncNew() func(conn net.Conn) cdc.Codec {
	return func(conn net.Conn) cdc.Codec {
		return New(conn)
	}
}

// Encoder 编码
func (c *codec) Encoder(data any) ([]byte, error) {
	// msgpack 编码
	buf, err := msgpack.Marshal(data)
	if err != nil {
		log.Println("Failed to encode:", err)
		return buf, err
	}
	return buf, nil
}

// Decoder 解码
func (c *codec) Decoder(b []byte) (any, error) {
	// msgpack 解码
	var data cdc.Response
	err := msgpack.Unmarshal(b, &data)
	if err != nil {
		log.Println("Failed to decode:", err)
		return &data, err
	}
	return &data, nil
}
