package zio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

var MSG_HEAD = []byte("@**@")

// Session 会话链接的结构体
type Session struct {
	conn net.Conn
	c    io.Closer
}

// NewSession 创建会话
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn, c: conn}
}

// 向连接中去写数据
func (s *Session) Write(data []byte) error {

	// 压缩数据
	//data, err := gzip.GzipEncode(data)
	//if err != nil {
	//	return err
	//}

	// 定义写数据的格式
	// 8字节头部 + 可变体的长度
	buf := make([]byte, 8+len(data))
	copy(buf[0:4], MSG_HEAD)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(data)))
	// 将整个数据，放到4后边
	copy(buf[8:], data)
	//log.Printf("bufbuf: %v", buf)
	_, err := s.conn.Write(buf)
	if err != nil {
		log.Printf("Write error: %v", err)
		return err
	}
	return nil
}

// 从连接中去读数据
func (s *Session) Read() ([]byte, error) {
	header := make([]byte, 8)
	//s.Conn.Read(io)
	// 按长度读取消息
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		if err != io.EOF {
			log.Printf("Read header error: %v", err)
		}
		return nil, err
	}

	if !bytes.Equal(header[0:4], MSG_HEAD) {
		err = errors.New("MSG_HEAD not found")
		log.Printf("%v", err)
		return nil, err
	}
	// 读取数据
	headerLen := header[4:8]
	dataLen := binary.BigEndian.Uint32(headerLen)
	body := make([]byte, dataLen)
	_, err = io.ReadFull(s.conn, body)

	// 解压缩数据
	//body, err := gzip.GzipDecode(body)
	//if err != nil {
	//	return nil, err
	//}

	if err != nil {
		log.Printf("Read body error: %v", err)
		return nil, err
	}
	return body, nil
}

func (s *Session) Close() error {
	log.Println("io close")
	return s.c.Close()
}
