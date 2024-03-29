package zrpc

import (
	"context"
	"errors"
	"fmt"
	"go/token"
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
	"zrpc/rpc/zio"
)

var e = (*error)(nil)
var typeOfInError = reflect.TypeOf(&e).Elem()
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// Ctx & Args & Reply 模式下的参数个数
var typeOfArgsNum = 5

// *Context 模式下的参数个数
var typeOfContextNum = 2

// Server 声明服务端
type Server struct {
	sd         center.ServeDiscovery
	addr       string
	serviceMap sync.Map

	codec func(conn net.Conn) codec.Codec
	io    func(conn net.Conn) zio.RWIo

	*RouterGroup
	pool      sync.Pool
	servePool sync.Pool
}

// Serve 服务
type Serve struct {
	codec      codec.Codec
	io         zio.RWIo
	serviceMap *sync.Map

	ServiceMethod string
	*RouterGroup
	pool *sync.Pool
}

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	argNum    int
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
		RouterGroup: &RouterGroup{
			handlers: nil,
		},
	}
	engine.pool.New = func() any {
		return &Serve{}
	}
	engine.servePool.New = func() any {
		return &Context{}
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
		var argType reflect.Type
		var replyType reflect.Type
		argNum := mtype.NumIn()
		if argNum == typeOfArgsNum {
			argType = mtype.In(2)
			if !isExportedOrBuiltinType(argType) {
				log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
				continue
			}
			replyType = mtype.In(3)
			if replyType.Kind() != reflect.Pointer {
				log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
				continue
			}
			if !isExportedOrBuiltinType(replyType) {
				log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
				continue
			}
			errorType := mtype.In(4)
			if errorType.Kind() != reflect.Pointer {
				log.Printf("rpc.Register: error type of method %q is not a pointer: %q\n", mname, errorType)
				continue
			}
			if errorType != typeOfInError {
				log.Printf("rpc.Register: error type of method %q is %q, must be error\n", mname, errorType)
				continue
			}
		} else {
			//log.Printf("rpc.Register: method %q has %d input parameters; needs exactly four\n", mname, mtype.NumIn())
			//continue
			replyType = mtype.In(1)
			if replyType.Kind() != reflect.Pointer {
				log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
				continue
			}
		}

		if mtype.NumOut() != 1 {
			log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			continue
		}
		if returnType := mtype.Out(0); returnType != typeOfError {
			log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			continue
		}
		methods[mname] = &methodType{method: method, argNum: argNum, ArgType: argType, ReplyType: replyType}

		// 注册方法集
		server.useHandle(name+"."+mname, func(c *Context) {
			method.Func.Call(c.inArgs)
		})
	}
	s.method = methods

	// 维护一个map
	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		log.Fatalf("rpc: service already defined: %s", sname)
	}

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

		//server.Serve(msgpack.New(conn), zio.NewSession(conn))
	}
}

func (serve *Serve) reset() {
	serve.ServiceMethod = ""
}

// Serve 建立服务
func (server *Server) Serve(codec codec.Codec, zio zio.RWIo) {
	serve := server.pool.Get().(*Serve)
	serve.reset()
	serve.codec = codec
	serve.io = zio
	serve.RouterGroup = server.RouterGroup
	serve.serviceMap = &server.serviceMap
	serve.pool = &server.servePool
	server.pool.Put(serve)
	go serve.ServeCodec()
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
func (serve *Serve) readRequest() (response *codec.Response, svc *service, mtype *methodType, keepReading bool, req bool, err error) {
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
	response = res.(*codec.Response)
	if err != nil {
		req = true
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		return
	}

	// 此次的coon是否已经结束
	keepReading = true

	// 根据读到的name，得到要调用的函数
	dot := strings.LastIndex(response.ServiceMethod, ".")
	if dot < 0 {
		err = errors.New(fmt.Sprintf("rpc: service/method request ill-formed: %s", response.ServiceMethod))
		return
	}

	// 存储path
	serve.ServiceMethod = response.ServiceMethod
	serviceName := response.ServiceMethod[:dot]
	methodName := response.ServiceMethod[dot+1:]

	// 查询注册的方法
	svci, ok := serve.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc: can't find service " + response.ServiceMethod)
		return
	}
	svc = svci.(*service)
	// 获取调用方法
	mtype = svc.method[methodName]

	return
}

// 结果返回客户端
func (serve *Serve) call(response *codec.Response, svc *service, mtype *methodType, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	// 捕获业务程序异常 防止崩溃
	defer func() {
		if err := recover(); err != nil {
			log.Printf("err.(error): %+v\n", err.(error))
			switch err.(type) {
			case string:
				serve.sendResponse(response, sending, errors.New(err.(string)))
			default:
				serve.sendResponse(response, sending, err.(error))
			}
			return
		}
	}()

	ctx := context.Background()
	// 包装参数
	inArgs := make([]reflect.Value, 0, typeOfArgsNum)
	if mtype.argNum == typeOfArgsNum {
		// 对象参数
		inArgs = append(inArgs, svc.rcvr)
		// ctx 参数
		inArgs = append(inArgs, reflect.ValueOf(ctx))
		// 获取Args对应的struct参数
		mapArgs := getArgsRValue(mtype, response.Args)
		inArgs = append(inArgs, mapArgs)
		// 返回值参数
		inArgs = append(inArgs, reflect.ValueOf(reflect.ValueOf(&response.Reply).Interface()))
		// 返回 error 参数
		inArgs = append(inArgs, reflect.ValueOf(&response.Error))
	}

	// 执行中间件
	h := serve.GetRoute(serve.ServiceMethod)
	if len(h) > 0 {
		c := serve.pool.Get().(*Context)
		c.reset()
		c.Ctx = ctx
		c.Args = response.Args
		c.Reply = &response.Reply
		c.Error = &response.Error
		c.inArgs = inArgs
		if mtype.argNum == typeOfContextNum {
			inContext := make([]reflect.Value, 0, typeOfContextNum)
			inContext = append(inContext, svc.rcvr)
			inContext = append(inContext, reflect.ValueOf(c))
			c.inArgs = inContext
		}
		c.handlers = h
		c.Next()
		serve.pool.Put(c)
	} else {
		serve.sendResponse(response, sending, errors.New("not find service"))
		return
	}

	// 返回Error类型
	errInter := response.Error
	var errReturn error
	if errInter != nil {
		errReturn = errInter.(error)
	} else {
		errReturn = (error)(nil)
	}

	// 发送客户端
	serve.sendResponse(response, sending, errReturn)
}

func (serve *Serve) sendResponse(response *codec.Response, sending *sync.Mutex, errReturn error) {
	sending.Lock()
	// 数据编码，返回给客户端
	respRPCData := codec.Response{ServiceMethod: response.ServiceMethod, Reply: response.Reply, Seq: response.Seq, Error: errReturn}
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
	}
	return v.Elem()
}
