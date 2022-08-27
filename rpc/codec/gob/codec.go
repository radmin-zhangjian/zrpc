package gob

import (
	"bytes"
	"encoding/gob"
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
	// gob 编码
	//得到字节数组的编码器
	var buf bytes.Buffer
	bufEnc := gob.NewEncoder(&buf)
	// 编码器对数据编码
	if err := bufEnc.Encode(data); err != nil {
		log.Printf("Encoder error: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decoder 解码
func (c *codec) Decoder(b []byte) (zio.Response, error) {
	// gob 解码
	buf := bytes.NewBuffer(b)
	// 得到字节数组解码器
	bufDec := gob.NewDecoder(buf)
	// 解码器对数据节码
	var data zio.Response
	if err := bufDec.Decode(&data); err != nil {
		log.Printf("Decoder error: %v", err)
		return data, err
	}
	return data, nil
}
