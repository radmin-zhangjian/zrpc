package main

import (
	"flag"
	"log"
	"net"
	"zrpc/example/service"
	v1 "zrpc/example/v1"
	v2 "zrpc/example/v2"
	"zrpc/rpc"
	"zrpc/rpc/codec/msgpack"
	"zrpc/rpc/socket/longSocket"
	"zrpc/rpc/zio"
)

// go run main.go -addr=127.0.0.1:8092
var (
	addr     = flag.String("addr", ":8092", "server address")
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

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
	srv := longSocket.NewServer(*addr, sd)

	// 将服务端方法，注册一下
	srv.RegisterName(new(service.Test), "service")
	srv.RegisterName(new(v1.Test), "v1")
	srv.RegisterName(new(v2.Test), "v2")

	// 启动服务
	lis := srv.Server()
	defer srv.Close(lis)

	// tpc连接map
	longSocket.ConnMap = make(map[string]*longSocket.Serve)

	// 开始监听
	accrpt(srv, lis)
}

func accrpt(srv *longSocket.Server, lis net.Listener) {
	// 开始监听
	for {
		conn, err := lis.Accept()
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			continue
		}

		session := zio.NewSession(conn)
		srv.Serve(msgpack.New(conn), session, conn)
	}
}
