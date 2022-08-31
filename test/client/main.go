package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
	"zrpc/example/service"
	v1 "zrpc/example/v1"
	v2 "zrpc/example/v2"
	"zrpc/rpc"
	"zrpc/rpc/center"
)

var (
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

// Args 参数
type Args struct {
	Id int64
	X  int64
	Y  int64
	Z  string
}

func closeCli() {
	cli = nil
}

// 自己定义数据格式的读写
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

	// 发现服务
	sd, err := rpc.ServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建客户端
	if cli == nil {
		cli, err = rpc.NewClient(sd, center.SelectMode(center.Random), true)
		defer closeCli()
		if err != nil {
			log.Fatal(err)
		}
	}

	// 压测
	startTime = GetCurrentTimeStampMS()
	wg := new(sync.WaitGroup)
	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		go Client(wg)
	}
	wg.Wait()
	dtime := GetCurrentTimeStampMS() - startTime
	fmt.Println("dTime:", dtime)

	// 同步rpc
	var reply any
	// 参数 struct 格式
	str := "我是rpc测试参数！！！"
	args := Args{
		Id: 2,
		X:  20,
		Z:  str,
	}
	errC := cli.Call("service.QueryUser", args, &reply)
	if errC != nil {
		fmt.Println("main.call.errC", errC)
	}
	reply1 := reply.(map[string]any)
	fmt.Println("main.call.reply", reply1["Age"])

	fmt.Println("==========================================")

	// 异步rpc
	var reply2 any
	call2 := cli.Go("v1.QueryInt", map[string]any{"Id": 10000, "msg": str}, &reply2, nil)
	<-call2.Done
	if call2.Error != nil {
		fmt.Printf("main.go.reply2.error: %v \n", call2.Error)
	}
	fmt.Printf("main.go.reply2: %v \n", reply2)

	time.Sleep(2 * time.Second)
}

var count int64 = 0
var startTime int64 = 0
var cli *rpc.Client

// 生成时间戳
func GetCurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func Client(wg *sync.WaitGroup) {
	defer wg.Done()
	var reply any
	// 参数 struct 格式
	args := Args{
		Id: 2,
		X:  20,
	}
	err := cli.Call("service.QueryUser", args, &reply)
	if err != nil {
		//fmt.Println("main.call.err", err)
	} else {
		//fmt.Println("main.call.reply", reply)
	}

	var reply2 any
	count++
	if count > 2 {
		count = 0
	}
	args2 := map[string]any{"Id": count}
	call := cli.Go("v1.QueryUser", args2, &reply2, nil)
	<-call.Done
	if call.Error != nil {
		//fmt.Printf("main.go.reply.error: %v \n", call.Error)
	} else {
		//fmt.Printf("main.go.reply: %v \n", reply2)
	}

}
