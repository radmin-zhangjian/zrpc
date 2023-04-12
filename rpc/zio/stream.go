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
	conn  net.Conn
	c     io.Closer
	token string
}

type RWIo interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
	SetToken(string)
	GetToken() string
}

// NewSession 创建会话
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn, c: conn, token: ""}
}

func FuncNew() func(conn net.Conn) RWIo {
	return func(conn net.Conn) RWIo {
		return NewSession(conn)
	}
}

func (s *Session) SetToken(token string) {
	s.token = token
}

func (s *Session) GetToken() string {
	return s.token
}

// 向连接中去写数据
func (s *Session) Write(data []byte) error {

	// 压缩数据
	//data, err := gzip.GzipEncode(data)
	//if err != nil {
	//	return err
	//}

	// 定义写数据的格式
	// 10字节头部 + 可变体的长度
	headerLen := 10
	if len(s.token) > 0 {
		headerLen += len(s.token)
	}

	buf := make([]byte, headerLen+len(data))
	copy(buf[0:4], MSG_HEAD)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(data)))

	// token
	if len(s.token) > 0 {
		binary.BigEndian.PutUint16(buf[8:10], uint16(len(s.token)))
		copy(buf[10:headerLen], []byte(s.token))
	} else {
		binary.BigEndian.PutUint16(buf[8:10], 0)
	}

	// 将整个数据，放到后边
	copy(buf[headerLen:], data)
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
	// 10字节头部
	headerLen := 10
	header := make([]byte, headerLen)
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
	bodyLen := header[4:8]
	dataLen := binary.BigEndian.Uint32(bodyLen)

	// auth
	authLen := header[8:10]
	authDataLen := binary.BigEndian.Uint16(authLen)
	if authDataLen > 0 {
		token := make([]byte, authDataLen)
		_, err = io.ReadFull(s.conn, token)
		s.token = string(token)
	}

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
