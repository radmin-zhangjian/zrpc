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
	"zrpc/rpc/codec/msgpack"
	"zrpc/rpc/zio"
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

	// protobuf 客户端
	//if cli == nil {
	//	//cli, err = rpc.NewClient(sd, center.SelectMode(center.Random), true)
	//	cli, err = rpc.LongClient(sd, center.SelectMode(center.Random))
	//	cli.SetOpt(pcd.New(cli.Conn)).SetOpt(zio.NewSession(cli.Conn))
	//	defer closeCli()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//
	//// 异步rpc
	//var reply3 any
	//str3 := "我是rpc测试参数222！！！"
	//args2 := &test.Args2{
	//	Id:    1,
	//	Param: str3,
	//}
	//inArgsAny, err := anypb.New(args2)
	//callproto := cli.Go("proto.GetUserList", inArgsAny, &reply3, nil)
	//<-callproto.Done
	//if callproto.Error != nil {
	//	fmt.Printf("proto.go.reply3.error: %v \n", callproto.Error)
	//} else {
	//	result := reply3.(*anypb.Any)
	//	unmarshal := &user.Reply{}
	//	anypb.UnmarshalTo(result, unmarshal, proto.UnmarshalOptions{})
	//	fmt.Printf("main.go.reply3: %v \n", unmarshal)
	//	fmt.Println("main.go.reply3.Detail", unmarshal.Detail)
	//	fmt.Println("main.go.reply3.Detail---", unmarshal.Detail["身高"])
	//	fmt.Println("main.go.reply3.Gender", unmarshal.Gender)
	//	fmt.Println("main.go.reply3.List", unmarshal.List)
	//	fmt.Println("main.go.reply3.List---", unmarshal.List[0].Model)
	//}
	//return

	// 创建客户端 (msgpack\json\gob)
	if cli == nil {
		//cli, err = rpc.NewClient(sd, center.SelectMode(center.Random), true)
		cli, err = rpc.LongClient(sd, center.SelectMode(center.Random))
		cli.SetOpt(msgpack.New(cli.Conn)).SetOpt(zio.NewSession(cli.Conn)).SetOptAuth("aaa111bbb222ccc3")
		defer closeCli()
		if err != nil {
			log.Fatal(err)
		}
	}

	// 压测
	startTime = GetCurrentTimeStampMS()
	wg := new(sync.WaitGroup)
	for i := 1; i <= 1000; i++ {
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
	} else {
		reply1 := reply.(map[string]any)
		fmt.Println("main.call.reply", reply1["Age"])
	}
	fmt.Println("main.call.reply", reply)

	fmt.Println("==========================================")

	// 异步rpc
	var reply2 any
	call2 := cli.Go("v1.QueryInt", map[string]any{"Id": 10000, "msg": str}, &reply2, nil)
	<-call2.Done
	if call2.Error != nil {
		fmt.Printf("main.go.reply2.error: %v \n", call2.Error)
	}
	fmt.Printf("main.go.reply2: %v \n", reply2)

	time.Sleep(1 * time.Second)
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
