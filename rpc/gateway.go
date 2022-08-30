package rpc

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"zrpc/rpc/center"
)

type Result struct {
	Code int64
	Data any
	Msg  any
}

type Http struct {
	addr string
}

// NewHttp 多语言网关
func NewHttp(addr string) *Http {
	return &Http{addr: addr}
}

// RegServe 初始化
func (s *Http) RegServe(sd center.ServeDiscovery, selectMode center.SelectAlgorithm, mode bool) *gin.Engine {
	// 启动模式
	gin.SetMode("debug")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", func(c *gin.Context) {
		result := Result{}
		var args map[string]any
		content := c.PostForm("content")
		if content == "" {
			result.Code = 1000
			result.Data = ""
			result.Msg = "content is empty"
			c.JSON(http.StatusOK, result)
			return
		}
		b := []byte(content)
		json.Unmarshal(b, &args)

		servicePath := c.PostForm("servicePath")
		serviceMethod := c.PostForm("serviceMethod")
		if servicePath == "" || serviceMethod == "" {
			result.Code = 1000
			result.Data = ""
			result.Msg = "servicePath or serviceMethod is empty"
			c.JSON(http.StatusOK, result)
			return
		}
		api := servicePath + "." + strings.ToUpper(serviceMethod[:1]) + serviceMethod[1:]

		client := NewClient(sd, selectMode, mode)
		var reply any
		call := client.Go(api, args, &reply, nil)
		<-call.Done
		if call.Error != nil {
			//fmt.Printf("main.go.reply.error: %v \n", call.Error)
			result.Code = 1000
			result.Data = ""
			result.Msg = call.Error.Error()
		} else {
			//fmt.Printf("main.go.reply: %v \n", reply)
			result.Code = 200
			result.Data = reply
			result.Msg = ""
		}
		c.JSON(http.StatusOK, result)
	})

	return router
}

// HttpServer 启动服务 & 优雅Shutdown（或重启）服务
func (s *Http) HttpServer(router *gin.Engine) {
	srv := &http.Server{
		Addr:         s.addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 5秒后优雅Shutdown服务
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt) //syscall.SIGKILL
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
	}
	log.Println("Server exiting")
}
