package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"zrpc/demo/gob_demo/header"
)

// Server 声明服务端
type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

// NewServer 构造方法
func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

// Register 服务端需要一个注册Register
// 第一个参数函数名，第二个传入真正的函数
func (s *Server) Register(method string, f any) {
	// 维护一个map
	// 若map已经有键了
	if _, ok := s.funcs[method]; ok {
		return
	}
	// 若map中没值，则将映射加入map，用于调用
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		panic(errors.New("invalid TypeOf Func"))
	}
	s.funcs[method] = fVal
}

// Run 服务端等待调用的方法
func (s *Server) Run() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("监听 %s err :%v", s.addr, err)
		return
	}
	for {
		conn, errL := lis.Accept()
		if errL != nil {
			continue
		}
		session := header.NewSession(conn)
		// 使用RPC方式读取数据
		b, errS := session.Read()
		if errS != nil {
			return
		}
		// 数据解码
		rpcData, errD := header.Decode(b)
		if errD != nil {
			return
		}
		// 根据读到的name，得到要调用的函数
		f, ok := s.funcs[rpcData.Name]
		if !ok {
			fmt.Printf("函数 %s 不存在", rpcData.Name)
			return
		}
		// 遍历解析客户端传来的参数,放切片里
		t := f.Type()
		if t.NumIn() != len(rpcData.Args) {
			log.Println("Err NumIn Lens Not Equal")
			return
		}
		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for i, arg := range rpcData.Args {
			v := reflect.ValueOf(arg)
			if t.In(i).Kind() != v.Type().Kind() {
				log.Println("Err parameter Type Mismatch,need:", t.In(i).Kind(), "but:", v.Type().Kind())
				return
			}
			inArgs = append(inArgs, v)
		}
		// 反射调用方法
		// 返回Value类型，用于给客户端传递返回结果,out是所有的返回结果
		out := f.Call(inArgs)
		// 遍历out ，用于返回给客户端，存到一个切片里
		outArgs := make([]interface{}, 0, len(out))
		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}
		// 数据编码，返回给客户端
		respRPCData := header.RpcData{rpcData.Name, outArgs}
		bytes, errE := header.Encode(respRPCData)
		if errE != nil {
			return
		}
		// 将服务端编码后的数据，写出到客户端
		errW := session.Write(bytes)
		if errW != nil {
			return
		}
	}
}
