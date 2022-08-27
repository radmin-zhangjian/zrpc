package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X, Y int
	S    string
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8091")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// 同步调用
	str := "qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]12345678909999"
	args := &Args{10000 * 1000, 20000 * 1000, str}
	var reply int
	err = client.Call("ServiceA.Add", args, &reply)
	if err != nil {
		log.Fatal("ServiceA.Add error:", err)
	}
	fmt.Printf("ServiceA.Add: %d+%d=%d\n", args.X, args.Y, reply)

	var reply3 int
	err = client.Call("ServiceB.Multiply", args, &reply3)
	if err != nil {
		log.Fatal("ServiceB.Multiply error:", err)
	}
	fmt.Printf("ServiceB.Multiply: %d+%d=%d\n", args.X, args.Y, reply3)

	// 异步调用
	var reply2 int
	divCall := client.Go("ServiceA.Add", args, &reply2, nil)
	replyCall := <-divCall.Done
	fmt.Println(replyCall.Error)
	fmt.Println(reply2)

}
