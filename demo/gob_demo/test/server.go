package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"zrpc/demo/gob_demo/header"
)

// 自己定义数据格式的读写
func main() {
	// 定义地址
	addr := "127.0.0.1:8000"
	my_data := "hello io i am god"
	// 等待组定义
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 写数据的协程
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		conn, _ := lis.Accept()
		session := header.NewSession(conn)
		err = session.Write([]byte(my_data))
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 读数据的协程
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		session := header.NewSession(conn)
		data, err := session.Read()
		if err != nil {
			log.Fatal(err)
		}
		if string(data) != my_data {
			log.Fatal(err)
		}
		fmt.Println(string(data))
	}()
	wg.Wait()
}
