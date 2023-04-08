package protobuf

import (
	"google.golang.org/protobuf/proto"
	"log"
	"net"
)

type codec struct {
	conn net.Conn
}

func New(conn net.Conn) *codec {
	return &codec{conn: conn}
}

// Encoder 编码
func (c *codec) Encoder(data any) ([]byte, error) {
	// 序列化结构体数据
	pdata := data.(*Response)
	buf, err := proto.Marshal(pdata)
	if err != nil {
		log.Fatalln("Failed to encode:", err)
		return buf, err
	}
	return buf, nil
}

// Decoder 解码
func (c *codec) Decoder(b []byte) (any, error) {
	//反序列化结构体
	data := &Response{}
	err := proto.Unmarshal(b, data)
	if err != nil {
		log.Fatalln("Failed to decode:", err)
		return data, err
	}
	return data, nil
}
