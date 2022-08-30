### protobuf 创建
protoc -I . --go_out=. --go-grpc_out=. ./response.proto

### server 服务端
````
package main

import (
	"flag"
	"log"
	"zrpc/rpc"
	"zrpc/rpc/proto/example"
	"zrpc/rpc/proto/rpcp"
)

// go run main.go -addr=127.0.0.1:8092
var (
	addr     = flag.String("addr", ":8092", "server address")
	registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
	basePath = flag.String("basepath", "/zrpc_center", "")
)

// 自己定义数据格式的读写
func main() {
	// 解析参数
	if !flag.Parsed() {
		flag.Parse()
	}

	// 创建服务发现
	sd, err := rpc.CreateServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}

	// 创建服务端
	srv := rpcp.NewServer(*addr, sd)
	// 将服务端方法，注册一下
	srv.RegisterName(new(example.Test), "service")
	// 启动服务
	srv.Serve()
}

````

### client 客户端
````
package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
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
	cli *rpcp.Client
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

	// 同步rpc
	var reply any
	// 参数 struct 格式
	str := "我是rpc测试参数！！！"
	args := &pd.Args{
		Id:    2,
		Param: str,
	}
	inArgsAny, err := ptypes.MarshalAny(args)
	errC := cli.Call("service.QueryProto", inArgsAny, &reply)
	if errC != nil {
		fmt.Println("main.call.errC", errC)
	} else {
		reply1 := reply.(*anypb.Any)
		unmarshal := &pd.Reply{}
		ptypes.UnmarshalAny(reply1, unmarshal)
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
	inArgsAny, err = ptypes.MarshalAny(args2)
	call2 := cli.Go("service.QueryProto2", inArgsAny, &reply2, nil)
	<-call2.Done
	if call2.Error != nil {
		fmt.Printf("main.go.reply2.error: %v \n", call2.Error)
	} else {
		result := reply2.(*anypb.Any)
		unmarshal := &pd.Reply{}
		ptypes.UnmarshalAny(result, unmarshal)
		fmt.Printf("main.go.reply2: %v \n", unmarshal)
	}

	time.Sleep(2 * time.Second)
}

````