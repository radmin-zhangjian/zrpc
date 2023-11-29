package service

import (
	"context"
	"log"
	"zrpc/example/common"
)

type Message struct {
}

// MessageArgsMap 参数 map传参时需要加上 json tag 区分大小写
type MessageArgsMap struct {
	Uid     int64  `json:"uid"`
	ToUid   int64  `json:"toUid"`
	Message string `json:"message"`
	Assign  string `json:"assign"`
}

func (m *Message) SendMap(ctx context.Context, arg MessageArgsMap, reply *any) error {
	log.Printf("Message.SendMap ID：%v, Assign: %v", arg.Uid, arg.Assign)
	code := 1000
	result := map[string]string{
		"message": arg.Message,
	}
	*reply = common.Response(code, common.GetMsg(code), result)
	return nil
}

// MessageArgs 参数
type MessageArgs struct {
	Uid     int64  `json:"uid"`
	ToUid   int64  `json:"toUid"`
	Message string `json:"message"`
	Assign  string `json:"assign"`
}

// Send 传递消息
func (m *Message) Send(ctx context.Context, arg MessageArgs, reply *any) error {
	log.Printf("Message.Send ID：%v, Assign: %v", arg.Uid, arg.Assign)
	code := 1000
	result := map[string]string{
		"message": arg.Message,
	}
	*reply = common.Response(code, common.GetMsg(code), result)
	return nil
}
