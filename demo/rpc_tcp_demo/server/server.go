package main

import (
	"log"
	"net"
	"net/rpc"
	"zrpc/demo/rpc_tcp_demo/server/services"
)

func main() {
	service := new(services.ServiceA)
	rpc.Register(service) // 注册RPC服务
	l, e := net.Listen("tcp", ":8091")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		conn, _ := l.Accept()
		rpc.ServeConn(conn)
	}
}
