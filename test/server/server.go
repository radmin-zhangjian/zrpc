package main

import (
	"encoding/gob"
	"flag"
	"log"
	"zrpc/example/service"
	v1 "zrpc/example/v1"
	v2 "zrpc/example/v2"
	"zrpc/rpc"
	"zrpc/rpc/codec/msgpack"
	"zrpc/rpc/proto/example"
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

	// gob 编解码时需要注册
	gob.Register(map[string]interface{}{})
	gob.Register(service.User{})
	gob.Register(v1.User{})
	gob.Register(v2.User{})

	// 创建服务发现
	sd, err := rpc.CreateServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建服务端
	srv := rpc.NewServer(*addr, sd)
	srv.SetOpt(msgpack.FuncNew())
	// 将服务端方法，注册一下
	//srv.Register(new(service.Test))
	srv.RegisterName(new(service.Test), "service")
	srv.RegisterName(new(v1.Test), "v1")
	srv.RegisterName(new(v2.Test), "v2")
	// protobuf 需要用proto前缀
	srv.RegisterProto(new(example.Test), "proto")
	// 启动服务
	srv.Accept(srv.Server())
}
