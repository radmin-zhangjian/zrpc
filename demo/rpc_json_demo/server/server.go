package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"zrpc/demo/rpc_json_demo/server/services"
)

func main() {
	service := new(services.ServiceA)
	rpc.Register(service) // 注册RPC服务
	ServiceB := new(services.ServiceB)
	rpc.Register(ServiceB) // 注册RPC服务
	l, e := net.Listen("tcp", ":8091")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		conn, err := l.Accept()
		//var buf []byte
		//newBuf := make([]byte, len(buf), 2*cap(buf)+10)
		//copy(newBuf, buf)
		//buf = newBuf
		//n, _ := conn.Read(buf[len(buf):cap(buf)])
		//buf = buf[0 : len(buf)+n]
		//log.Println(buf)
		//io.ServeCodec(jsonrpc.NewServerCodec(conn))
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			fmt.Println("new client")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}
