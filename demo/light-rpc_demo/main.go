package main

import (
	"fmt"
	"log"
	"time"
	"zrpc/demo/light-rpc_demo/lrpc"
)

var count = 0
var startTime int64 = 0
var client *lrpc.RPCClient

// 生成时间戳
func GetCurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / 1e6
}

type Test struct {
}

func (this *Test) Add(a, b int, c string) int {
	log.Println(c)
	return a + b
}

func (this *Test) Hello() {
	fmt.Println("hello!")
}

func (this *Test) GetName() (string, int) {
	return "Peter", 20
}

func (this *Test) GetTest(a, b int, c string) string {
	return c
}

func main() {
	go Server()

	startTime = GetCurrentTimeStampMS()
	//go Client()

	if client == nil {
		client = lrpc.NewRPCClient("localhost", "2048")
		client.Dial()
	}
	str := "qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]12345678909999MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM0000"
	arg := []interface{}{
		1,
		2,
		str,
	}
	client.CallReply("Test.GetTest", arg, func(arg ...interface{}) {
		log.Printf("Test.GetTest: %v", arg)
	})

	for {
		a := 1
		a++
		time.Sleep(time.Second)
	}
}

func Client() {
	if client == nil {
		client = lrpc.NewRPCClient("localhost", "2048")
		client.Dial()
	}

	str := "qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]1234567890zxcvbnm,.';lkjhgffdsaqwertyuiop[]=-0987654321qwerfcvbnmjuiopl" +
		"qwertyuioplkjhgfdsazxcvbnm,.[]12345678909999MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM"
	arg := []interface{}{
		1,
		2,
		str,
	}
	client.CallReply("Test.Add", arg, func(arg ...interface{}) {
		count++
		if count < 1 {
			Client()
		} else {
			dtime := GetCurrentTimeStampMS() - startTime
			fmt.Println("dTime:", dtime)

			client.Call("Test.Hello", nil)

			client.CallReply("Test.GetName", nil, func(arg ...interface{}) {
				fmt.Println(arg)
			})
		}
	})
}

func Server() {
	server := lrpc.NewRPCServer("2048")
	server.Register(&Test{})
	server.Accept()
}
