package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	addr = flag.String("addr", ":10000", "addr server")
)

type Http struct {
	addr string
}

// NewHttp 多语言网关
func NewHttp(addr string) *Http {
	return &Http{addr: addr}
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
	}
	log.Println("Server exiting")
}

// RegServe 初始化
func (s *Http) RegServe() *gin.Engine {
	// 启动模式
	gin.SetMode("debug")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/get", get)
	router.GET("/set", set)
	router.POST("/chanSet", chanSet)

	return router
}

//获取ip
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

//获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

var chanMap = make(map[string]chan any)

func get(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	key := c.Query("key")

	statusChan := make(chan any)
	chanMap[key] = statusChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(20)*time.Second)
	defer cancel()

	var value any
	select {
	case <-ctx.Done():
		value = "nil"
	case value = <-statusChan:
	}

	var result = make(map[string]any)
	result["key"] = key
	result["value"] = value

	c.JSON(http.StatusOK, result)
}

func set(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	IP, _ := externalIP()
	queryCustomerUrl := "http://" + IP.String() + ":10000/chanSet"
	param := url.Values{}
	param.Add("key", key)
	param.Add("value", value)
	resp, _ := http.PostForm(queryCustomerUrl, param)
	log.Println("resp: ", resp)

	// 返回
	var result = make(map[string]any)
	result["key"] = key
	result["value"] = value

	c.JSON(http.StatusOK, result)
}

func chanSet(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")

	mu := new(sync.Mutex)
	mu.Lock()
	chanMap[key] <- value
	mu.Unlock()
}

func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}
	hp := NewHttp(*addr)
	// 注册服务
	router := hp.RegServe()
	// 启动http服务
	hp.HttpServer(router)
	// 可以在启动一个http服务 区分get和set请求
	//hp.HttpServer(router)
}
