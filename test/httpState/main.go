package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sync"
	"time"
	"zrpc/rpc/center/redis"
	"zrpc/test/httpState/server"
	"zrpc/utils"
)

var (
	addr    = flag.String("addr", ":10000", "addr server")
	chanMap = make(map[string]chan any)
	prefix  = "state:"
	timeOut = 20
)

func Response(code string, msg string, data any) any {
	result := make(map[string]any)
	result["code"] = code
	result["msg"] = msg
	result["data"] = data
	return result
}

func get(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	var result = make(map[string]any)

	key := c.Query("key")

	statusChan := make(chan any)
	mu := new(sync.Mutex)
	mu.Lock()
	chanMap[key] = statusChan
	mu.Unlock()

	defer func() {
		if _, ok := chanMap[key]; ok {
			mu.Lock()
			delete(chanMap, key)
			mu.Unlock()
		}
		close(statusChan)
		//RedisDel(key)
	}()

	IP, err := utils.ExternalIP()
	if err != nil {
		c.JSON(http.StatusOK, Response("1000", err.Error(), nil))
		return
	}
	rKey := prefix + key
	err = redis.SetStringExp(rKey, IP.String(), 60*time.Second)
	if err != nil {
		c.JSON(http.StatusOK, Response("1000", err.Error(), nil))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		result["value"] = "nil"
	case value := <-statusChan:
		result["key"] = key
		result["value"] = value
	}

	c.JSON(http.StatusOK, Response("200", "success", result))
}

func set(c *gin.Context) {
	key := c.Query("key")
	value := c.Query("value")

	// 返回
	var result = make(map[string]any)
	result["data"] = nil

	rKey := prefix + key
	IP, err := redis.GetString(rKey)
	if err != nil {
		c.JSON(http.StatusOK, Response("1000", err.Error(), result))
		return
	}

	queryCustomerUrl := "http://" + IP.(string) + ":10000/chanSet"
	param := url.Values{}
	param.Add("key", key)
	param.Add("value", value)
	http.PostForm(queryCustomerUrl, param)
	//resp, _ := http.PostForm(queryCustomerUrl, param)
	//log.Println("resp: ", resp)

	result["key"] = key
	result["value"] = value

	c.JSON(http.StatusOK, Response("200", "success", result))
}

func chanSet(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	mu := new(sync.Mutex)
	mu.Lock()
	if ch, ok := chanMap[key]; ok {
		ch <- value
	}
	mu.Unlock()

	c.JSON(http.StatusOK, Response("200", "success", nil))
}

func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}
	hp := server.NewHttp(*addr)
	redis.InitRedis("127.0.0.1", "6379", "", 0, 100)
	// 注册服务
	gin.SetMode("debug")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("test/httpState/templates/*")
	router.GET("/httpState.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "httpState.html", gin.H{})
	})
	router.GET("/get", get)
	router.GET("/set", set)
	router.POST("/chanSet", chanSet)
	// 启动http服务
	hp.HttpServer(router)
	// 可以在启动一个http服务 区分get和set请求
	//hp.HttpServer(router)
}
