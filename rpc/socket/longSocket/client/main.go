package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"zrpc/rpc"
	"zrpc/rpc/center"
	"zrpc/rpc/socket/longSocket"
)

var (
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
	cli      *longSocket.Client
)

func closeCli() {
	cli = nil
}

// 自己定义数据格式的读写
func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}

	// 发现服务
	sd, err := rpc.ServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建客户端
	if cli == nil {
		cli, err = longSocket.NewClient(sd, center.SelectMode(center.Random))
		defer closeCli()
		if err != nil {
			log.Fatal(err)
		}
	}

	// 输出
	go Call(longSocket.DoneCall)

	// 异步rpc
	var reply any
	// map形式传参数时 需要注意 大小写 这块容易出错
	//args := map[string]any{
	//	"uid":     22,
	//	"toUid":   1,
	//	"message": "我是message测试参数！！！",
	//	"assign":  "single",
	//}
	//cli.Go("message.SendMap", args, &reply, nil)

	// struct 形式传参数时 接收方
	type MessageArgs struct {
		Uid     int    `json:"uid"`
		ToUid   int    `json:"toUid"`
		Message string `json:"message"`
		Assign  string `json:"assign"`
	}
	messageArgs := MessageArgs{
		Uid:     22,
		ToUid:   1,
		Message: "你好，我是章鱼",
		Assign:  "single",
	}
	cli.Go("message.Send", messageArgs, &reply, nil)

	for {

	}
}

func Call(call *longSocket.Done) {
	for {
		select {
		case c := <-call.Done:
			if c.Error != nil {
				fmt.Printf("main.go.Done.error: %v \n", c.Error)
			}
			data, _ := json.Marshal(c.Reply)
			fmt.Printf("main.go.Done: %v \n", string(data))
		}
	}
}
