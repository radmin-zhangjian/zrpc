# zrpc
***
说明: zrpc实现了以下基本功能  
1、自定义协议字节流  
2、定义编解码，支持多种方法  
3、注册服务方法  
4、反射调用方法及参数，根据struct方式和map方式传参  
5、支持同步和异步调用  
6、服务发现：redis实现(根据路径查询和redis-scan方式)、其他方式待实现  
7、负载均衡：随机和轮询、其他方式待实现  
8、捕获业务程序异常防止崩溃  
9、支持多语言 网关(http)  
10、php通过socket调用zrpc  

### 服务端
go run test/server/main.go  -addr=127.0.0.1:8091  
go run test/server/main.go  -addr=127.0.0.1:8092  
go run test/server/main.go  -addr=127.0.0.1:8093  

### 客户端
go run test/client/main.go  

### php客户端 (socket方式)
test/client/client.php  

### 网关 支持多语言
go run test/gateway/gateway.go  

### php 连接网关 (curl方式)
test/gateway/gateway.php  

### 快速开始 - 服务端
```
package main

import (
    "encoding/gob"
    "flag"
    "log"
    "zrpc/example/service"
    v1 "zrpc/example/v1"
    v2 "zrpc/example/v2"
    "zrpc/rpc"
)

// go run main.go -addr=127.0.0.1:8092
var (
    addr     = flag.String("addr", ":8092", "server address")
    registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
    basePath = flag.String("basepath", "/zrpc_center", "")
)

// 自己定义数据格式的读写
func main() {
    // 解析参数
    if !flag.Parsed() {
    flag.Parse()
    }

	// 创建服务发现
	sd, err := rpc.CreateServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}
 
	// 创建服务端
	srv := rpc.NewServer(*addr, sd)
	// 将服务端方法，注册一下
	//srv.Register(new(service.Test))
	srv.RegisterName(new(service.Test), "service")
	srv.RegisterName(new(v1.Test), "v1")
	srv.RegisterName(new(v2.Test), "v2")
	// 启动服务
	srv.Serve()
}
```

### 快速开始 - 客户端
````
package main

import (
    "encoding/gob"
    "flag"
    "fmt"
    "log"
    "sync"
    "time"
    "zrpc/example/service"
    v1 "zrpc/example/v1"
    v2 "zrpc/example/v2"
    "zrpc/rpc"
    "zrpc/rpc/center"
)

var (
    registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
    basePath = flag.String("basepath", "/zrpc_center", "")
    cli *rpc.Client
)

// Args 参数
type Args struct {
    Id int64
    X  int64
    Y  int64
    Z  string
}

// 自己定义数据格式的读写
func main() {
    // 解析参数
    if !flag.Parsed() {
    flag.Parse()
    }

	// 发现服务
	sd, err := rpc.ServiceDiscovery(*basePath, *registry, "", 0, 100)
	if err != nil {
		log.Fatal(err)
	}
 
	// 创建客户端
	if cli == nil {
		cli = rpc.NewClient(sd, center.SelectMode(center.Random), true)
		defer closeCli()
	}
 
	// 同步rpc
	var reply any
	// 参数 struct 格式
	str := "我是rpc测试参数！！！"
	args := Args{
		Id: 2,
		X:  20,
		Z:  str,
	}
	errC := cli.Call("service.QueryUser", args, &reply)
	if errC != nil {
		fmt.Println("main.call.errC", errC)
	}
	reply1 := reply.(map[string]any)
	fmt.Println("main.call.reply", reply1["Age"])
 
	fmt.Println("==========================================")
 
	// 异步rpc
	var reply2 any
	call2 := cli.Go("v1.QueryInt", map[string]any{"Id": 10000, "msg": str}, &reply2, nil)
	<-call2.Done
	if call2.Error != nil {
		fmt.Printf("main.go.reply2.error: %v \n", call2.Error)
	}
	fmt.Printf("main.go.reply2: %v \n", reply2)
 
	time.Sleep(2 * time.Second)
}
````


### 快速开始 - php客户端
````
(new TestController())->actionTest()

class TestController extends Controller
{
const HEAD_MSG = "@**@";
public $socket;

    // 析构函数
    public function __destruct() {
        socket_close($this->socket);
    }
 
    public function newRpc($host, $port) {
        $this->socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP) or die("Unable to create socket");
        @socket_connect($this->socket, $host, $port) or die("Connect error.");
        if ($err = socket_last_error($this->socket)) {
            socket_close($this->socket);
            die(socket_strerror($err));
        }
    }
 
    public function build($method, $data) {
        $response = [
            "ServiceMethod" => $method,
            "Args" => $data,
            "Reply" => null,
        ];
        return $response;
    }
 
    public function write($response) {
        $response = msgpack_pack($response);
        $len = strlen($response);
        $buf = pack("a4", self::HEAD_MSG);
        $buf .= pack("N", $len);
        $buf .= pack("a".$len, $response);
        return socket_write ($this->socket , $buf, strlen($buf));
    }
 
    public function read() {
        // 读取数据
        $hm = socket_read($this->socket, 4, PHP_BINARY_READ);
        if ($hm != self::HEAD_MSG) {
            print_r("HEAD_MSG ERROR: " . $hm);
        }
        $hLen = socket_read($this->socket, 4, PHP_BINARY_READ);
        $hLen = unpack("N", $hLen);
        $hLen = $hLen[1];
        $buffer = socket_read($this->socket, $hLen, PHP_BINARY_READ);
        $buffer = msgpack_unpack($buffer);
        return $buffer;
    }
 
    public function call($api, $data) {
        // 打包数据
        $buf = $this->build($api, $data);
 
        // 发送数据
        $this->write($buf);
 
        // 读取数据
        $buffer = $this->read();
        return $buffer;
    }
 
    // PHP 调用 zrpc 服务
    public function actionTest()
    {
        // 连接服务
        $this->newRpc("127.0.0.1", "8092");
 
        // 数据
        $data["Id"] = 1;
        $data["X"] = 20;
        $data["Z"] = "aaasssdddfffggghhh";
        $data["msg"] = "msg000";
 
        // 接口
        $api[] = "service.QueryUser";
        $api[] = "v1.QueryInt";
 
        // 调用
        for ($i=0; $i<10000; $i++) {
            $rand = mt_rand(0, 1);
            $buffer = $this->call($api[$rand], $data);
            var_dump(json_encode($buffer));
        }
 
    }
}
````
### 快速开始 - 网关
1、启动网关
````
package main

import (
    "flag"
    "log"
    "zrpc/rpc"
)

var (
    addr     = flag.String("addr", "127.0.0.1:8060", "addr server")
    registry = flag.String("registry", "redis://127.0.0.1:6379", "registry address")
    basePath = flag.String("basepath", "/zrpc_center", "")
)

func main() {
    // 解析参数
    if !flag.Parsed() {
        flag.Parse()
    }
    // http new
    http := rpc.NewHttp(*addr)
    // 发现服务
    sd, err := rpc.ServiceDiscovery(*basePath, *registry, "", 0, 100)
    if err != nil {
        log.Fatal(err)
    }
    // 注册服务
    router := http.RegServe(sd)
    // 启动http服务
    http.HttpServer(router)
}
````
2、客户端连接网关 php例子
````
    $url = 'http://localhost:8060/';
    $data['servicePath'] = 'v1';
    $data['serviceMethod'] = 'queryInt';
    $data['content'] = '{"Id": 10000, "msg": "测试api"}';
 
    $this->curl = new CurlRequest($url);
    $this->curl->setPostFeilds($data);
    $data = $this->curl->post();
    $data = json_decode($data, true);
    var_dump($data);
 
    // curl
    class CurlRequest
    {
        public $version = "1.0";
 
        public $handler;
 
        public $timeOut = 30;
 
        public $header = ['Expect:'];
 
        public $referer;
 
        public $postFields;
 
        public $url;
 
        public $ssl = false;
 
        public $type = "get";
 
        public $agent = "see curl/1.0";
 
        public $returnData;
 
        public function __construct($url)
        {
            $this->url = $url = trim($url);
            $this->handler = curl_init();
            curl_setopt($this->handler, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_0);
            curl_setopt($this->handler, CURLOPT_RETURNTRANSFER, true);
            curl_setopt($this->handler, CURLOPT_URL, $url);
            //支持https
            $this->ssl = stripos($url, 'https://') === 0 ? true : false;
            if ($this->ssl) {
                curl_setopt($this->handler, CURLOPT_SSL_VERIFYPEER, false);
                curl_setopt($this->handler, CURLOPT_SSL_VERIFYHOST, false);
            }
        }
 
        public function setTimeOut($timeOut){
            $this->timeOut = $timeOut;
        }
 
        public function setOpt($opt, $value)
        {
            curl_setopt($this->handler, $opt, $value);
        }
 
        public function setHeader($header)
        {
            $this->header = array_merge($this->header, $header);
        }
 
        public function setCookie($cookie)
        {
            curl_setopt($this->handler, CURLOPT_COOKIE, $cookie);
        }
 
        public function setReferer($referer)
        {
            $this->referer = $referer;
        }
 
        public function setPostFeilds($postFields)
        {
            $this->postFields = $postFields;
        }
 
        public function setAgent($agent){
            $this->agent = $agent;
        }
 
        public function exec()
        {
            if (!empty($this->referer)) {
                curl_setopt($this->handler, CURLOPT_REFERER, $this->referer);
            }else{
                curl_setopt($this->handler, CURLOPT_AUTOREFERER, true);
            }
 
            //set header
            $this->setHeader(['PHP-SEE-TID:' . 1]);
            $this->setHeader(['PHP-SEE-SEQ:' . 2]);
 
            curl_setopt($this->handler,CURLOPT_HTTPHEADER,$this->header);
            curl_setopt($this->handler, CURLOPT_USERAGENT,$this->agent);
            curl_setopt($this->handler, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_0);
            curl_setopt($this->handler, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4);
            curl_setopt($this->handler, CURLOPT_CONNECTTIMEOUT, $this->timeOut);
            curl_setopt($this->handler, CURLOPT_TIMEOUT, $this->timeOut);
 
            if($this->type=='post'){
                curl_setopt($this->handler, CURLOPT_POST, true);
                curl_setopt($this->handler, CURLOPT_POSTFIELDS, $this->postFields);
            }
 
            $this->returnData = curl_exec($this->handler);
            $httpcode = curl_getinfo($this->handler, CURLINFO_HTTP_CODE);
            if ($errorNo = curl_errno($this->handler) || $httpcode != 200) {
                //error message
                $errorMsg = curl_error($this->handler);
                \See::$log->warning("curl error, url:%s, type:%s, postData:%s, errorNo:%s, errorMsg:%s, httpcode:%s, return:%s",$this->url,$this->type,$this->postFields,$errorNo,$errorMsg,$httpcode,$this->returnData);
            }else{
                \See::$log->trace("curl success, url:%s, type:%s, postData:%s, httpcode:%s, return:%s",$this->url,$this->type,$this->postFields,$httpcode,$this->returnData);
            }
            curl_close($this->handler);
            return $this->returnData;
        }
 
        public function get(){
            $this->type = 'get';
            return $this->exec();
        }
 
        public function post(){
            $this->type ='post';
            return $this->exec();
        }
 
    }
````