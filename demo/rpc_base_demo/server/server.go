package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"zrpc/demo/rpc_base_demo/server/services"
)

func main() {
	service := new(services.ServiceA)
	rpc.Register(service) // 注册RPC服务
	rpc.HandleHTTP()      // 基于HTTP协议
	l, e := net.Listen("tcp", ":8091")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
