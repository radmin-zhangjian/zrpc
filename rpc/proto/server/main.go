package main

import (
	"flag"
	"log"
	"zrpc/rpc"
	"zrpc/rpc/proto/example"
	"zrpc/rpc/proto/rpcp"
)

// go run main.go -addr=127.0.0.1:8092
var (
	addr     = flag.String("addr", ":8092", "server address")
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

// 自己定义数据格式的读写
func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}

	// 创建服务发现
	sd, err := rpc.CreateServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建服务端
	srv := rpcp.NewServer(*addr, sd)
	// 将服务端方法，注册一下
	srv.RegisterName(new(example.Test), "proto")
	// 启动服务
	srv.Serve()
}
