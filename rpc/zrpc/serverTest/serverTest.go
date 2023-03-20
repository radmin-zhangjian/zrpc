package main

import (
	"encoding/gob"
	"flag"
	"log"
	"zrpc/rpc"
	"zrpc/rpc/zrpc"
	"zrpc/rpc/zrpc/example/service"
	v1 "zrpc/rpc/zrpc/example/v1"
	v2 "zrpc/rpc/zrpc/example/v2"
)

// go run main.go -addr=127.0.0.1:8092
var (
	addr     = flag.String("addr", ":8092", "server address")
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

// +++ 增加了中间件的服务 +++
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

	// 创建中间件
	//r := rpc.NewRouterGroup()
	//r.Use(func(c *rpc.Context) {
	//	log.Println("/Use")
	//	c.Next()
	//	log.Println("/Use next")
	//})
	//r.UseHandle("/v1/test", func(c *rpc.Context) {
	//	log.Println("/v1/test111")
	//	c.Next()
	//	log.Println("/v1/test111-222")
	//}, func(c *rpc.Context) {
	//	log.Println("/v1/test222")
	//})
	//r.UseHandle("/v2/test", func(c *rpc.Context) {
	//	log.Println("/v2/test")
	//})
	//h := r.GetRoute("/v1/test")
	//log.Printf("handles: %+v\n", h)
	//context := rpc.NewContext()
	//context.Test(h)

	// 注册服务
	sd, err := rpc.CreateServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建服务端
	srv := zrpc.NewServer(*addr, sd)
	// 将服务端方法，注册一下
	srv.RegisterName(new(service.Test), "service")
	srv.RegisterName(new(v1.Test), "v1")
	srv.RegisterName(new(v2.Test), "v2")
	// 中间件 实例
	srv.Use(func(c *zrpc.Context) {
		log.Println("srv use ==========")
		c.Next()
		log.Println("srv use next ==========")
	})
	srv.UseHandle("v1.QueryInt",
		func(c *zrpc.Context) {
			log.Println("v1.QueryInt ++++++++++")
			log.Println("c.Args ++++++++++", c.Args)
		}, func(c *zrpc.Context) {
			log.Println("v1.QueryInt ----------")
			c.Next()
			log.Println("c.QueryInt next ----------")
		})
	// 启动服务
	srv.Accept(srv.Server())
}
