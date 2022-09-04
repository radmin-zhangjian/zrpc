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
	"sync"
	"time"
	"zrpc/rpc/center"
)

var client *Client
var ErrShutdownnum int64 = 0

type Result struct {
	Code int64 `json:"code"`
	Data any   `json:"data"`
	Msg  any   `json:"msg"`
}

type Http struct {
	addr string
}

// NewHttp 多语言网关
func NewHttp(addr string) *Http {
	return &Http{addr: addr}
}

func send(c *gin.Context, sd center.ServeDiscovery, selectMode center.SelectAlgorithm, mode bool, wg *sync.WaitGroup) {
	//defer wg.Done()
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

	if mode == false {
		// 短连接模式
		nc, err := ShortClient(sd, selectMode)
		if err != nil {
			if err == ErrDiscovery {
				result.Code = 1000
				result.Data = ""
				result.Msg = err.Error()
				c.JSON(http.StatusOK, result)
			} else {
				time.Sleep(2 * time.Second)
				send(c, sd, selectMode, mode, wg)
			}
			return
		}
		client = nc
	} else {
		// 长连接模式
		if client == nil {
			nc, err := LongClient(sd, selectMode)
			if err != nil {
				if err == ErrDiscovery {
					result.Code = 1000
					result.Data = ""
					result.Msg = err.Error()
					c.JSON(http.StatusOK, result)
				} else {
					time.Sleep(2 * time.Second)
					send(c, sd, selectMode, mode, wg)
				}
				return
			}
			client = nc
		}
	}

	var reply any
	call := client.Go(api, args, &reply, nil)
	<-call.Done
	if call.Error != nil {
		errs := call.Error
		if errs == ErrShutdown {
			ErrShutdownnum++
			if ErrShutdownnum < 10 {
				log.Println("ErrShutdown: ", errs)
				client = nil
				//wg.Add(1)
				send(c, sd, selectMode, mode, wg)
				return
			}
			return
		}
		result.Code = 1000
		result.Data = ""
		result.Msg = call.Error.Error()
	} else {
		ErrShutdownnum = 0
		result.Code = 200
		result.Data = reply
		result.Msg = ""
	}
	c.JSON(http.StatusOK, result)
}

// RegServe 初始化
func (s *Http) RegServe(sd center.ServeDiscovery, selectMode center.SelectAlgorithm, mode bool) *gin.Engine {
	// 启动模式
	gin.SetMode("debug")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", func(c *gin.Context) {
		wg := &sync.WaitGroup{}
		//wg.Add(1)
		//go send(c, sd, selectMode, mode, wg)
		//wg.Wait()
		send(c, sd, selectMode, mode, wg)
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
