package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	pprofGin "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "net/http/pprof"
	"zrpc/rpc/center/redis"
	"zrpc/test/workerPool/api"
	"zrpc/test/workerPool/server"
	"zrpc/utils/workerPool"
)

var (
	addr       = flag.String("addr", ":10001", "addr server")
	cpuprofile = flag.String("cpuprofile", "fabonacci.prof", "write cpu profile to file")
)

func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}

	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.StartCPUProfile(f)
	//	defer pprof.StopCPUProfile()
	//}

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

	//性能调优监视
	//authStr := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(config.ProducerUsername+":"+config.ProducerPassword)))
	authStr := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("zhang:123456")))
	pprofGroup := router.Group("/producer", func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth != authStr {
			//c.Header("www-Authenticate", "Basic")
			//c.AbortWithStatus(http.StatusUnauthorized)
			//return
		}
		c.Next()
	})
	pprofGin.RouteRegister(pprofGroup, "worker_pool_pprof")

	redis.InitRedis("127.0.0.1", "6379", "", 0, 100)
	// 启动http服务
	hp.HttpServer(router)
}
