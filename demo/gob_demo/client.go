package main

import (
	"errors"
	"net"
	"reflect"
	"zrpc/demo/gob_demo/header"
)

// Client 声明客户端
type Client struct {
	conn net.Conn
}

// NewClient 构造方法
func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

// 实现通用的RPC客户端
// 传入访问的函数名
// fPtr指向的是函数原型
//var select fun xx(User)
//cli.callRPC("selectUser",&select)
func (c *Client) call(method string, fPtr any) {
	// 通过反射，获取fPtr未初始化的函数原型
	fn := reflect.ValueOf(fPtr).Elem()
	if fn.Type().Kind() != reflect.Func {
		panic(errors.New("err fPtr Not Func"))
	}
	// 需要另一个函数，作用是对第一个函数参数操作
	f := func(args []reflect.Value) []reflect.Value {
		// 处理参数
		inArgs := make([]any, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		// 连接
		session := header.NewSession(c.conn)
		// 编码数据
		reqData := header.RpcData{Name: method, Args: inArgs}
		b, err := header.Encode(reqData)
		if err != nil {
			panic(err)
		}
		// 写数据
		err = session.Write(b)
		if err != nil {
			panic(err)
		}

		// 服务端发过来返回值，此时应该读取和解析
		respBytes, errR := session.Read()
		if errR != nil {
			panic(errR)
		}
		// 解码
		respData, errD := header.Decode(respBytes)
		if errD != nil {
			panic(errR)
		}
		// 处理服务端返回的数据
		outArgs := make([]reflect.Value, 0, len(respData.Args))
		for i, arg := range respData.Args {
			// 必须进行nil转换
			if arg == nil {
				// reflect.Zero()会返回类型的零值的value
				// .out()会返回函数输出的参数类型
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(arg))
		}
		return outArgs
	}
	// 完成原型到函数调用的内部转换
	// 参数1是reflect.Type
	// 参数2 f是函数类型，是对于参数1 fn函数的操作
	// fn是定义，f是具体操作
	v := reflect.MakeFunc(fn.Type(), f)
	// 为函数fPtr赋值，过程
	fn.Set(v)
}
