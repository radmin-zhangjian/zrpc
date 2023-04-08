package rpcp

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
	"strings"
	"sync"
	"zrpc/rpc/center"
	pcd "zrpc/rpc/codec/protobuf"
	"zrpc/rpc/zio"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfAnypd = reflect.TypeOf((*anypb.Any)(nil)).Elem()

//type ServerCodec interface {
//	Encoder(pcd.Response) ([]byte, error)
//	Decoder(b []byte) (pcd.Response, error)
//}
type ServerCodec interface {
	Encoder(any) ([]byte, error)
	Decoder(b []byte) (any, error)
}

type ServerIo interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}

// Server 声明服务端
type Server struct {
	codec      ServerCodec
	sd         center.ServeDiscovery
	addr       string
	serviceMap sync.Map
	io         ServerIo
}

// Serve 服务
type Serve struct {
	codec      ServerCodec
	io         ServerIo
	serviceMap sync.Map
}

type methodType struct {
	method  reflect.Method
	ArgType reflect.Type
}

type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value
	method map[string]*methodType
}

// NewServer 构造方法
func NewServer(addr string, sd center.ServeDiscovery) *Server {
	return &Server{addr: addr, sd: sd}
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

// Serve 监听tcp
func (server *Server) Serve() (lis net.Listener) {
	lis, e := net.Listen("tcp", server.addr)
	if e != nil {
		log.Fatalf("监听 %s err :%v", server.addr, e)
		server.Close(lis)
		return
	}
	defer server.Close(lis)

	// 注册服务中心
	server.sd.ServeRegister(server.addr)
	// 心跳监测
	go server.sd.TickerHeartbeatTask(server.addr)

	server.Accept(lis)
	return
}

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			if conn != nil {
				conn.Close()
			}
			continue
		}

		server.ServerConn(pcd.New(conn), zio.NewSession(conn))
	}
}

func (server *Server) ServerConn(codec ServerCodec, zio ServerIo) {
	serve := &Serve{codec: codec, io: zio}
	serve.serviceMap = server.serviceMap
	go serve.ServeCodec()
}

// ServeCodec 调用接口
func (serve *Serve) ServeCodec() {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		response, svc, mtype, keepReading, req, err := serve.readRequest()
		if err != "" {
			if !keepReading {
				break
			}
			if !req {
				serve.sendResponse(response, sending, err)
			}
			continue
		}
		wg.Add(1)
		go serve.call(response, svc, mtype, sending, wg)
	}
	wg.Wait()
	serve.io.Close()
}

// 读取并解析参数
func (serve *Serve) readRequest() (response *pcd.Response, svc *service, mtype *methodType, keepReading bool, req bool, err string) {
	// 使用RPC方式读取数据
	b, errR := serve.io.Read()
	if errR != nil {
		err = errR.Error()
		req = true
		if errR == io.EOF || errR == io.ErrUnexpectedEOF {
			return
		}
		return
	}

	// 数据解码
	res, errC := serve.codec.Decoder(b)
	response = res.(*pcd.Response)
	if errC != nil {
		err = errR.Error()
		req = true
		if errC == io.EOF || errC == io.ErrUnexpectedEOF {
			return
		}
		return
	}

	// 此次的coon是否已经结束
	keepReading = true

	// 根据读到的name，得到要调用的函数
	dot := strings.LastIndex(response.ServiceMethod, ".")
	if dot < 0 {
		err = errors.New(fmt.Sprintf("rpc: service/method request ill-formed: %s", response.ServiceMethod)).Error()
		return
	}
	serviceName := response.ServiceMethod[:dot]
	methodName := response.ServiceMethod[dot+1:]
	// 查询注册的方法
	//service, _ := server.serviceMap.Load(serviceName)
	//rValue := reflect.ValueOf(service)
	//rType := reflect.TypeOf(service)
	// 获取调用方法
	//method, _ = rType.MethodByName(methodName)

	// 查询注册的方法
	svci, ok := serve.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc: can't find service " + response.ServiceMethod).Error()
		return
	}
	svc = svci.(*service)
	// 获取调用方法
	mtype = svc.method[methodName]

	return
}

// 结果返回客户端
func (serve *Serve) call(response *pcd.Response, svc *service, mtype *methodType, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	// 捕获业务程序异常 防止崩溃
	defer func() {
		if err := recover(); err != nil {
			e := err.(error)
			serve.sendResponse(response, sending, e.Error())
			return
		}
	}()

	// 包装参数
	inArgs := make([]reflect.Value, 0, 4)
	// 对象参数
	inArgs = append(inArgs, svc.rcvr)
	// ctx 参数
	inArgs = append(inArgs, reflect.ValueOf(context.Background()))
	// 获取Args对应的struct参数
	mapArgs := getArgsRValue(mtype, response.Args)
	inArgs = append(inArgs, mapArgs)

	// 获取方法
	method := mtype.method
	// 反射调用方法
	returnValues := method.Func.Call(inArgs)

	// 返回Error类型
	reply := returnValues[0].Interface()
	errInter := returnValues[1].Interface()
	var errReturn string
	if errInter != nil {
		errReturn = errInter.(error).Error()
	} else {
		errReturn = ""
	}
	response.Reply = reply.(*anypb.Any)
	response.Error = errReturn
	// 发送客户端
	serve.sendResponse(response, sending, errReturn)
}

func (serve *Serve) sendResponse(response *pcd.Response, sending *sync.Mutex, errReturn string) {
	sending.Lock()
	// 数据编码，返回给客户端
	respRPCData := pcd.Response{ServiceMethod: response.ServiceMethod, Reply: response.Reply, Seq: response.Seq, Error: errReturn}
	bytes, errE := serve.codec.Encoder(&respRPCData)
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
	return reflect.ValueOf(reflect.ValueOf(args).Interface())
}
