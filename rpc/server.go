package rpc

import (
	"context"
	"errors"
	"fmt"
	"go/token"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"zrpc/rpc/center"
	"zrpc/rpc/codec"
	"zrpc/rpc/codec/msgpack"
	pcd "zrpc/rpc/codec/protobuf"
	"zrpc/rpc/zio"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// Server 声明服务端
type Server struct {
	sd   center.ServeDiscovery
	addr string
	pool sync.Pool // serve pool

	codec      func(conn net.Conn) codec.Codec
	io         func(conn net.Conn) zio.RWIo
	serviceMap sync.Map
	authFunc   func(string) bool
}

// Serve 服务
type Serve struct {
	codec      codec.Codec
	io         zio.RWIo
	serviceMap *sync.Map
	authFunc   func(string) bool
}

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value
	method map[string]*methodType
}

// NewServer 构造方法
func NewServer(addr string, sd center.ServeDiscovery) *Server {
	engine := &Server{
		addr: addr,
		sd:   sd,
	}
	engine.pool.New = func() any {
		return &Serve{}
	}
	return engine
}

// SetOpt 自定义设置opt
func (server *Server) SetOpt(opt any) *Server {
	switch opt.(type) {
	case func(conn net.Conn) codec.Codec:
		server.codec = opt.(func(conn net.Conn) codec.Codec)
	case func(conn net.Conn) zio.RWIo:
		server.io = opt.(func(conn net.Conn) zio.RWIo)
	}
	return server
}

func (server *Server) SetOptAuth(f func(string) bool) *Server {
	server.authFunc = f
	return server
}

// Register 服务端注册服务
func (server *Server) Register(rcvr any) error {
	return server.register(rcvr, "")
}

// RegisterName 服务端注册服务
func (server *Server) RegisterName(rcvr any, name string) error {
	return server.register(rcvr, name)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

func (server *Server) register(rcvr any, name string) error {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	// 判断是否是指针
	if s.rcvr.Kind() != reflect.Pointer {
		return errors.New("invalid TypeOf Struct")
	}
	s.typ = reflect.TypeOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if name != "" {
		sname = name
	}
	if sname == "" {
		e := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(e)
		return errors.New(e)
	}
	// 判断是否是可调用方法
	if !token.IsExported(sname) && name == "" {
		e := "rpc.Register: type " + sname + " is not exported"
		log.Print(e)
		return errors.New(e)
	}
	s.name = sname
	log.Println(sname)
	// 注册服务
	methods := make(map[string]*methodType)
	fTyp := s.typ
	for i := 0; i < fTyp.NumMethod(); i++ {
		method := fTyp.Method(i)
		mtype := method.Type
		mname := method.Name
		if !method.IsExported() {
			continue
		}
		if mtype.NumIn() != 4 {
			log.Printf("rpc.Register: method %q has %d input parameters; needs exactly four\n", mname, mtype.NumIn())
			continue
		}
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			continue
		}
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Pointer {
			log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			continue
		}
		if !isExportedOrBuiltinType(replyType) {
			log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			continue
		}
		if mtype.NumOut() != 1 {
			log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			continue
		}
		if returnType := mtype.Out(0); returnType != typeOfError {
			log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	s.method = methods

	// 维护一个map
	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		log.Fatalf("rpc: service already defined: %s", sname)
	}

	// 维护一个map
	//if _, dup := server.serviceMap.LoadOrStore(sname, rcvr); dup {
	//	log.Fatalf("rpc: service already defined: %s", sname)
	//}

	return nil
}

func (server *Server) RegisterProto(rcvr any, name string) error {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	// 判断是否是指针
	if s.rcvr.Kind() != reflect.Pointer {
		return errors.New("invalid TypeOf Struct")
	}
	s.typ = reflect.TypeOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if name != "" {
		sname = name
	}
	if sname == "" {
		e := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(e)
		return errors.New(e)
	}
	// 判断是否是可调用方法
	if !token.IsExported(sname) && name == "" {
		e := "rpc.Register: type " + sname + " is not exported"
		log.Print(e)
		return errors.New(e)
	}
	s.name = sname
	log.Println(sname)
	// 注册服务
	methods := make(map[string]*methodType)
	fTyp := s.typ
	for i := 0; i < fTyp.NumMethod(); i++ {
		method := fTyp.Method(i)
		mtype := method.Type
		mname := method.Name
		if !method.IsExported() {
			continue
		}
		if mtype.NumIn() != 3 {
			log.Printf("rpc.Register: method %q has %d input parameters; needs exactly four\n", mname, mtype.NumIn())
			continue
		}
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			continue
		}
		if mtype.NumOut() != 2 {
			log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			continue
		}
		//log.Println(mtype.Out(0))
		//if returnType := mtype.Out(0); returnType != typeOfAnypd {
		//	log.Printf("rpc.Register: return type of method %q is %q, must be *anypb.Any\n", mname, returnType)
		//	continue
		//}
		if returnType := mtype.Out(1); returnType != typeOfError {
			log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType}
	}
	s.method = methods

	// 维护一个map
	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		log.Fatalf("rpc: service already defined: %s", sname)
	}

	// 维护一个map
	//if _, dup := server.serviceMap.LoadOrStore(sname, rcvr); dup {
	//	log.Fatalf("rpc: service already defined: %s", sname)
	//}

	return nil
}

func (server *Server) Close(lis net.Listener) {
	server.sd.Close(server.addr)
	lis.Close()
}

// Server 启动服务
func (server *Server) Server() (lis net.Listener) {
	lis, e := net.Listen("tcp", server.addr)
	if e != nil {
		log.Fatalf("监听 %s err :%v", server.addr, e)
		server.Close(lis)
		return
	}

	// 注册服务中心
	server.sd.ServeRegister(server.addr)
	// 心跳监测
	go server.sd.TickerHeartbeatTask(server.addr)

	return
}

// Accept 监听tcp
func (server *Server) Accept(lis net.Listener) {
	defer server.Close(lis)
	for {
		conn, err := lis.Accept()
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			continue
		}

		var cc codec.Codec
		if server.codec == nil {
			cc = msgpack.New(conn)
		} else {
			cc = server.codec(conn)
		}

		var rwio zio.RWIo
		if server.io == nil {
			rwio = zio.NewSession(conn)
		} else {
			rwio = server.io(conn)
		}

		server.Serve(cc, rwio)
	}
}

// Serve 建立服务
func (server *Server) Serve(codec codec.Codec, zio zio.RWIo) {
	serve := server.pool.Get().(*Serve)
	serve.reset()
	serve.codec = codec
	serve.io = zio
	serve.serviceMap = &server.serviceMap
	serve.authFunc = server.authFunc
	server.pool.Put(serve)
	go serve.ServeCodec()
}

func (serve *Serve) reset() {

}

// ServeCodec 调用接口
func (serve *Serve) ServeCodec() {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		response, svc, mtype, keepReading, req, err := serve.readRequest()
		if err != nil {
			//if err == io.EOF {
			//	break
			//}
			if err.Error() == "auth error" {
				serve.sendResponse(response, sending, err)
				break
			}
			if !keepReading {
				break
			}
			if !req {
				serve.sendResponse(response, sending, err)
			}
			continue
		}
		wg.Add(1)
		serve.call(response, svc, mtype, sending, wg)
	}
	wg.Wait()
	serve.io.Close()
}

// 读取并解析参数
func (serve *Serve) readRequest() (response any, svc *service, mtype *methodType, keepReading bool, req bool, err error) {
	// 使用RPC方式读取数据
	b, err := serve.io.Read()
	if err != nil {
		req = true
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		return
	}

	// 数据解码
	res, err := serve.codec.Decoder(b)
	//response = res.(*codec.Response)
	response = res
	if err != nil {
		req = true
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		return
	}

	// auth认证
	if serve.authFunc != nil {
		af := serve.authFunc(serve.io.GetToken())
		if !af {
			err = errors.New("auth error")
			log.Printf("%s", "auth error")
			return
		}
	}

	// 此次的coon是否已经结束
	keepReading = true

	// 获取服务方法
	switch res.(type) {
	case *pcd.Response:
		smethod := res.(*pcd.Response)
		// 根据读到的name，得到要调用的函数
		dot := strings.LastIndex(smethod.ServiceMethod, ".")
		if dot < 0 {
			err = errors.New(fmt.Sprintf("rpc: service/method request ill-formed: %s", smethod.ServiceMethod))
			return
		}
		serviceName := smethod.ServiceMethod[:dot]
		methodName := smethod.ServiceMethod[dot+1:]

		// 查询注册的方法
		svci, ok := serve.serviceMap.Load(serviceName)
		if !ok {
			err = errors.New("rpc: can't find service " + smethod.ServiceMethod)
			return
		}
		svc = svci.(*service)
		// 获取调用方法
		mtype = svc.method[methodName]
	default:
		smethod := res.(*codec.Response)
		// 根据读到的name，得到要调用的函数
		dot := strings.LastIndex(smethod.ServiceMethod, ".")
		if dot < 0 {
			err = errors.New(fmt.Sprintf("rpc: service/method request ill-formed: %s", smethod.ServiceMethod))
			return
		}
		serviceName := smethod.ServiceMethod[:dot]
		methodName := smethod.ServiceMethod[dot+1:]
		// 查询注册的方法
		//service, _ := server.serviceMap.Load(serviceName)
		//rValue := reflect.ValueOf(service)
		//rType := reflect.TypeOf(service)
		// 获取调用方法
		//method, _ = rType.MethodByName(methodName)

		// 查询注册的方法
		svci, ok := serve.serviceMap.Load(serviceName)
		if !ok {
			err = errors.New("rpc: can't find service " + smethod.ServiceMethod)
			return
		}
		svc = svci.(*service)
		// 获取调用方法
		mtype = svc.method[methodName]
	}

	return
}

// 结果返回客户端
func (serve *Serve) call(response any, svc *service, mtype *methodType, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	// 捕获业务程序异常 防止崩溃
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case string:
				serve.sendResponse(response, sending, errors.New(err.(string)))
			default:
				serve.sendResponse(response, sending, err.(error))
			}
			return
		}
	}()

	switch response.(type) {
	case *pcd.Response:
		resp := response.(*pcd.Response)
		// 包装参数
		inArgs := make([]reflect.Value, 0, 4)
		// 对象参数
		inArgs = append(inArgs, svc.rcvr)
		// ctx 参数
		inArgs = append(inArgs, reflect.ValueOf(context.Background()))
		// 获取Args对应的struct参数
		mapArgs := getArgsRValue(mtype, resp.Args)
		inArgs = append(inArgs, mapArgs)

		// 获取方法
		method := mtype.method
		// 反射调用方法
		returnValues := method.Func.Call(inArgs)

		// 返回Error类型
		reply := returnValues[0].Interface()
		errInter := returnValues[1].Interface()
		var errReturn error
		if errInter != nil {
			errReturn = errInter.(error)
		} else {
			errReturn = nil
		}
		resp.Reply = reply.(*anypb.Any)
		if errReturn == nil {
			resp.Error = ""
		} else {
			resp.Error = errReturn.Error()
		}

		// 发送客户端
		serve.sendResponse(resp, sending, errReturn)

	default:
		resp := response.(*codec.Response)
		// 包装参数
		inArgs := make([]reflect.Value, 0, 4)
		// 对象参数
		inArgs = append(inArgs, svc.rcvr)
		// ctx 参数
		inArgs = append(inArgs, reflect.ValueOf(context.Background()))
		// 获取Args对应的struct参数
		mapArgs := getArgsRValue(mtype, resp.Args)
		inArgs = append(inArgs, mapArgs)
		// 返回值参数
		inArgs = append(inArgs, reflect.ValueOf(reflect.ValueOf(&resp.Reply).Interface()))

		// 获取方法
		method := mtype.method
		// 反射调用方法
		returnValues := method.Func.Call(inArgs)
		// 返回Error类型
		errInter := returnValues[0].Interface()
		var errReturn error
		if errInter != nil {
			errReturn = errInter.(error)
		} else {
			errReturn = nil
		}

		// 发送客户端
		serve.sendResponse(resp, sending, errReturn)
	}

	// 发送客户端
	//serve.sendResponse(response, sending, errReturn)
}

func (serve *Serve) sendResponse(response any, sending *sync.Mutex, errReturn error) {
	sending.Lock()
	// 数据编码，返回给客户端
	var reqData any
	switch response.(type) {
	case *pcd.Response:
		resp := response.(*pcd.Response)
		reqData = pcd.Response{ServiceMethod: resp.ServiceMethod, Reply: resp.Reply, Seq: resp.Seq, Error: resp.Error}
	default:
		resp := response.(*codec.Response)
		reqData = codec.Response{ServiceMethod: resp.ServiceMethod, Reply: resp.Reply, Seq: resp.Seq, Error: errReturn}
	}
	respRPCData := reqData

	//respRPCData := codec.Response{ServiceMethod: response.ServiceMethod, Reply: response.Reply, Seq: response.Seq, Error: errReturn}
	bytes, errE := serve.codec.Encoder(respRPCData)
	if errE != nil {
		//return
	}

	// 将服务端编码后的数据，写出到客户端
	err := serve.io.Write(bytes)
	if err != nil {
		return
	}
	sending.Unlock()
}

// 包装参数给方法
func getArgsRValue(mtype *methodType, args any) reflect.Value {
	switch args.(type) {
	case *anypb.Any:
		return reflect.ValueOf(reflect.ValueOf(args).Interface())
	default:
		// 将interface{}类型的map转换为  map[string]interface{}
		argMap := args.(map[string]any)

		// 查找Args参数  reflect.Type类型 *xxx.Args.Elem()
		//t := method.Type.In(2)
		t := mtype.ArgType
		// 获取结构体Args对象  reflect.Value类型
		v := reflect.New(t)
		// 判断是否是struct类型
		if t.Kind() == reflect.Struct {
			argNum := t.NumField()
			for i := 0; i < argNum; i++ {
				// 获取指定字段的反射值
				f := v.Elem().Field(i)
				// 获取struct的指定字段
				stf := t.Field(i)
				// 获取tag
				name := strings.Split(stf.Tag.Get("json"), ",")[0]
				// 判断是否为忽略字段
				if name == "-" {
					continue
				}
				// 判断是否为空，若为空则使用字段本身的名称获取value值
				if name == "" {
					name = stf.Name
				}
				// 获取value值
				val, ok := argMap[name]
				if !ok {
					continue
				}
				// 获取指定字段的类型
				kind := v.Elem().Field(i).Kind()
				// 若字段为指针类型
				if kind == reflect.Ptr {
					// 获取对应字段的kind
					kind = f.Type().Elem().Kind()
				}
				// 赋值
				switch kind {
				case reflect.Int:
				case reflect.Int64:
					res, _ := strconv.ParseInt(fmt.Sprint(val), 10, 64)
					v.Elem().Field(i).SetInt(res)
				case reflect.String:
					v.Elem().Field(i).SetString(fmt.Sprint(val))
				}
			}
		} else if t.Kind() == reflect.Map {
			return reflect.ValueOf(reflect.ValueOf(args).Interface())
		} else {
			return reflect.ValueOf(reflect.ValueOf(args).Interface())
		}
		return v.Elem()
	}
}
