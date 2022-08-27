package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

// 自己定义数据格式的读写
func main() {
	// 编码中有一个字段是interface{}时，要注册一下
	gob.Register(User{})
	addr := "127.0.0.1:8093"
	// 创建服务端
	srv := NewServer(addr)
	// 将服务端方法，注册一下
	srv.Register("queryUser", queryUser)
	// 服务端等待调用
	go srv.Run()

	// 客户端
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("err")
	}
	// 创建客户端对象
	cli := NewClient(conn)
	// 需要声明函数原型
	var query func(int) (User, error)
	cli.call("queryUser", &query)
	// 得到查询结果
	u, errQ := query(1)
	if errQ != nil {
		fmt.Println("err")
	}
	fmt.Println(u)
}

// 定义用户对象
type User struct {
	Name string
	Age  int
}

// 用于测试用户查询的方法
func queryUser(uid int) (User, error) {
	user := make(map[int]User)
	// 假数据
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	user[2] = User{"ww", 22}
	// 模拟查询用户
	if u, ok := user[uid]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("%d err", uid)
}
