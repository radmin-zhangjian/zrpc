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
func (c *codec) Encoder(data Response) ([]byte, error) {
	// 序列化user结构体数据
	buf, err := proto.Marshal(&data)
	if err != nil {
		log.Fatalln("Failed to encode:", err)
		return buf, err
	}
	return buf, nil

}

// Decoder 解码
func (c *codec) Decoder(b []byte) (Response, error) {
	//反序列化user结构体
	data := Response{}
	err := proto.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Failed to decode:", err)
		return data, err
	}
	return data, nil
}
