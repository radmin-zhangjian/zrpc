package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"zrpc/rpc/center/redis"
	"zrpc/test/httpState/server"
	"zrpc/test/workerPool/api"
	"zrpc/utils/workerPool"
)

var (
	addr = flag.String("addr", ":10001", "addr server")
)

func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}

	// 开启工作池
	dispatcher := workerPool.NewDispatcher(workerPool.MaxWorker)
	dispatcher.Run()
	//log.Println("当前协程数：", runtime.NumGoroutine())

	// http
	hp := server.NewHttp(*addr)

	// 注册服务
	gin.SetMode("debug")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/get", api.To)
	redis.InitRedis("127.0.0.1", "6379", "", 0, 100)
	// 启动http服务
	hp.HttpServer(router)
}
