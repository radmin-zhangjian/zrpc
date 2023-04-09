package zrpc

import (
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	"zrpc/rpc/center"
	"zrpc/rpc/codec"
	"zrpc/rpc/codec/msgpack"
	"zrpc/rpc/zio"
)

var debugLog = false
var ErrShutdown = errors.New("connection is shut down")
var ErrDiscovery = errors.New("service not found")

// Call 返回调用方
type Call struct {
	// 访问的函数
	ServiceMethod string
	// 访问时的参数
	Args any
	// 返回数据
	Reply any
	// 错误
	Error error
	// call
	Done chan *Call
}

// Client 声明客户端
type Client struct {
	codec codec.Codec
	io    zio.RWIo
	Conn  net.Conn

	selectMode string

	seq      uint64
	pending  map[uint64]*Call
	mutex    sync.Mutex
	shutdown bool
}

// ClientConn 构造方法
func ClientConn(sd center.ServeDiscovery, sm center.SelectAlgorithm) (net.Conn, error) {
	// 发现服务
	disArrs, err := sd.ServeDiscovery()
	if err != nil {
		return nil, ErrDiscovery
	}

	// 路由/负载均衡
	addr := sm.Algorithm(disArrs)

	// 失败处理
	//failMode   string // 失败处理

	// 链接服务端
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewClient 构造方法
func NewClient(conn net.Conn, codec codec.Codec, zio zio.RWIo, mode bool) *Client {
	cli := &Client{io: zio, codec: codec, Conn: conn, pending: make(map[uint64]*Call)}
	if mode == true {
		go cli.input()
	} else {
		go cli.inputNoCycle()
	}
	return cli
}

// ClientRpc 构造方法
func ClientRpc(sd center.ServeDiscovery, sm center.SelectAlgorithm, mode bool) (*Client, error) {
	// 发现服务
	disArrs, err := sd.ServeDiscovery()
	if err != nil {
		return nil, ErrDiscovery
	}

	// 路由/负载均衡
	addr := sm.Algorithm(disArrs)

	// 失败处理
	//failMode   string // 失败处理

	// 链接服务端
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// 创建客户端对象
	cli := NewClient(conn, msgpack.New(conn), zio.NewSession(conn), mode)
	//cli := NewClient(conn, nil, nil, mode)
	return cli, nil
}

func ShortClient(sd center.ServeDiscovery, sm center.SelectAlgorithm) (*Client, error) {
	return ClientRpc(sd, sm, false)
}

func LongClient(sd center.ServeDiscovery, sm center.SelectAlgorithm) (*Client, error) {
	return ClientRpc(sd, sm, true)
}

// SetOpt 自定义设置opt
func (c *Client) SetOpt(opt any) *Client {
	switch opt.(type) {
	case codec.Codec:
		c.codec = opt.(codec.Codec)
	case zio.RWIo:
		c.io = opt.(zio.RWIo)
	}
	return c
}

// Call 同步RPC客户端
func (c *Client) Call(serviceMethod string, args any, reply any) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

// Go 异步RPC客户端
func (c *Client) Go(serviceMethod string, args any, reply any, done chan *Call) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	if done == nil {
		done = make(chan *Call, 10)
	}
	call.Done = done
	c.send(call)
	return call
}

// 发送服务端数据
func (c *Client) send(call *Call) {
	// 存储 call
	c.mutex.Lock()
	if c.shutdown {
		c.mutex.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := c.seq
	c.seq++
	c.pending[seq] = call
	c.mutex.Unlock()

	// 处理参数
	var inArgs any
	mapArgs := make(map[string]any)
	argsKind := reflect.ValueOf(call.Args).Kind()
	if argsKind == reflect.Struct {
		v := reflect.ValueOf(call.Args)
		t := reflect.TypeOf(call.Args)
		argNum := v.NumField()
		c.mutex.Lock()
		for i := 0; i < argNum; i++ {
			mapArgs[t.Field(i).Name] = v.Field(i).Interface()
		}
		c.mutex.Unlock()
		inArgs = mapArgs
	} else if argsKind == reflect.Map {
		inArgs = call.Args
	}

	// 编码数据
	reqData := codec.Response{ServiceMethod: call.ServiceMethod, Args: inArgs, Seq: seq}
	b, err := c.codec.Encoder(reqData)
	if err != nil {
		log.Printf("rpc encode: %v", err)
	}
	// 写数据
	err = c.io.Write(b)
	if err != nil {
		log.Printf("rpc write: %v", err)
	}
}

// 处理服务端返回的数据
func (c *Client) input() {
	var err error
	for err == nil {
		// 服务端发过来返回值，读取和解析
		respBytes, errR := c.io.Read()
		if errR != nil {
			err = errors.New("reading error body: " + errR.Error())
		}
		// 解码
		res, errD := c.codec.Decoder(respBytes)
		response := res.(*codec.Response)
		if errD != nil {
			err = errors.New("reading error body: " + errD.Error())
			break
		}

		// 获取返回 call
		seq := response.Seq
		c.mutex.Lock()
		call := c.pending[seq]
		delete(c.pending, seq)
		c.mutex.Unlock()

		// 处理返回数据
		switch {
		case call == nil:

		case response.Error != nil:
			call.Error = response.Error
			call.done()
		default:
			// 处理服务端返回的数据
			replay := call.Reply.(*any)
			*replay = response.Reply
			call.done()
		}
	}

	c.mutex.Lock()
	c.shutdown = true
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
	c.mutex.Unlock()
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		if debugLog {
			log.Println("rpc: discarding Call reply due to insufficient Done chan capacity")
		}
	}
}

// 处理服务端返回的数据
func (c *Client) inputNoCycle() {
	var err error
	// 服务端发过来返回值，读取和解析
	respBytes, errR := c.io.Read()
	if errR != nil {
		err = errors.New("reading error body: " + errR.Error())
	}
	// 解码
	res, errD := c.codec.Decoder(respBytes)
	response := res.(*codec.Response)
	if errD != nil {
		err = errors.New("reading error body: " + errD.Error())
	}

	// 获取返回 call
	seq := response.Seq
	c.mutex.Lock()
	call := c.pending[seq]
	delete(c.pending, seq)
	c.mutex.Unlock()

	// 处理返回数据
	switch {
	case call == nil:

	case response.Error != nil:
		call.Error = response.Error
		call.done()
	default:
		// 处理服务端返回的数据
		//var outArgs []byte
		//for _, arg := range response.Args {
		//	for _, a := range arg.([]byte) {
		//		outArgs = append(outArgs, a)
		//	}
		//}
		//fmt.Println("response.outArgs: ", string(outArgs))

		replay := call.Reply.(*any)
		*replay = response.Reply
		//call.Reply = replay
		//log.Printf("replay------: %#v", replay)
		call.done()
	}

	c.mutex.Lock()
	c.shutdown = true
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
	c.mutex.Unlock()
	c.Conn.Close()
}
