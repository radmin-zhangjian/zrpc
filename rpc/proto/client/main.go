package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"sync"
	"time"
	"zrpc/rpc"
	"zrpc/rpc/center"
	pd "zrpc/rpc/proto"
	"zrpc/rpc/proto/rpcp"
)

var (
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

func closeCli() {
	cli = nil
}

// 自己定义数据格式的读写
func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}

	// 发现服务
	sd, err := rpc.ServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建客户端
	if cli == nil {
		cli, err = rpcp.NewClient(sd, center.SelectMode(center.Random), true)
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
	args := &pd.Args{
		Id:    2,
		Param: str,
	}
	//inArgsAny, err := ptypes.MarshalAny(args)
	inArgsAny, err := anypb.New(args)
	errC := cli.Call("service.QueryProto", inArgsAny, &reply)
	if errC != nil {
		fmt.Println("main.call.errC", errC)
	} else {
		reply1 := reply.(*anypb.Any)
		unmarshal := &pd.Reply{}
		//ptypes.UnmarshalAny(reply1, unmarshal)
		anypb.UnmarshalTo(reply1, unmarshal, proto.UnmarshalOptions{})
		fmt.Println("main.call.reply", unmarshal)
	}

	fmt.Println("==========================================")

	// 异步rpc
	var reply2 any
	str = "我是rpc测试参数222！！！"
	args2 := &pd.Args2{
		Id:    1,
		Param: str,
	}
	inArgsAny, err = anypb.New(args2)
	call2 := cli.Go("service.QueryProto2", inArgsAny, &reply2, nil)
	<-call2.Done
	if call2.Error != nil {
		fmt.Printf("main.go.reply2.error: %v \n", call2.Error)
	} else {
		result := reply2.(*anypb.Any)
		unmarshal := &pd.Reply{}
		anypb.UnmarshalTo(result, unmarshal, proto.UnmarshalOptions{})
		fmt.Printf("main.go.reply2: %v \n", unmarshal)
	}

	time.Sleep(2 * time.Second)
}

var count int64 = 0
var startTime int64 = 0
var cli *rpcp.Client

// 生成时间戳
func GetCurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func Client(wg *sync.WaitGroup) {
	defer wg.Done()
	var reply any
	// 参数 struct 格式
	args := &pd.Args{
		Id:    2,
		Param: "msg",
	}
	inArgsAny, _ := anypb.New(args)
	err := cli.Call("service.QueryProto", inArgsAny, &reply)
	if err != nil {
		//fmt.Println("main.call.err", err)
	} else {
		//fmt.Println("main.call.reply", reply)
	}

	var reply2 any
	call := cli.Go("service.QueryProto", inArgsAny, &reply2, nil)
	<-call.Done
	if call.Error != nil {
		//fmt.Printf("main.go.reply.error: %v \n", call.Error)
	} else {
		//fmt.Printf("main.go.reply: %v \n", reply2)
	}

}
