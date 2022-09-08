package main

import (
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
	str := "我是rpc测试参数！！！"
	cli.Go("v1.QueryInt", map[string]any{"Id": 10000, "msg": str}, &reply, nil)

	cli.Go("v1.QueryUser", map[string]any{"Id": 2, "msg": str}, &reply, nil)

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
			fmt.Printf("main.go.Done: %v \n", c.Reply)
		}
	}
}
