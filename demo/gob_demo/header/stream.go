package header

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"net"
)

// Session 会话连接的结构体
type Session struct {
	conn net.Conn
}

// NewSession 创建会话
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn}
}

// 向连接中去写数据
func (s *Session) Write(data []byte) error {
	// 定义写数据的格式
	// 4字节头部 + 可变体的长度
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	// 将整个数据，放到4后边
	copy(buf[4:], data)
	_, err := s.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

// 从连接中去读数据
func (s *Session) Read() ([]byte, error) {
	header := make([]byte, 4)
	//s.Conn.Read(io)
	// 按长度读取消息
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		return nil, err
	}
	// 读取数据
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// RpcData 定义RPC交互的数据结构
type RpcData struct {
	// 访问的函数
	Name string
	// 访问时的参数
	Args []interface{}
}

// Encode 编码
func Encode(data RpcData) ([]byte, error) {
	//得到字节数组的编码器
	var buf bytes.Buffer
	bufEnc := gob.NewEncoder(&buf)
	// 编码器对数据编码
	if err := bufEnc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode 解码
func Decode(b []byte) (RpcData, error) {
	buf := bytes.NewBuffer(b)
	// 得到字节数组解码器
	bufDec := gob.NewDecoder(buf)
	// 解码器对数据节码
	var data RpcData
	if err := bufDec.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}
