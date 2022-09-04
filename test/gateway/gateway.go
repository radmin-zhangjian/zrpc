package main

import (
	"flag"
	"log"
	"zrpc/rpc"
	"zrpc/rpc/center"
)

var (
	addr     = flag.String("addr", "127.0.0.1:8060", "addr server")
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}
	// http new
	http := rpc.NewHttp(*addr)
	// 发现服务
	sd, err := rpc.ServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}
	// 注册服务  mode=false 短连接模式  mode=true 长连接模式
	router := http.RegServe(sd, center.SelectMode(center.RoundRobin), false)
	// 启动http服务
	http.HttpServer(router)
}
