<?php

namespace app\console;


use see\console\Controller;

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
        // 数据头
        $buf = pack("a4", self::HEAD_MSG);
        // 数据体长度
        $buf .= pack("N", $len);

        // token 长度
        $token_len = strlen(self::TOKEN);
        $buf .= pack("n", $token_len);
        // token 主体
        $buf .= pack("a".$token_len, self::TOKEN);

        // 数据主体
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

        // token
        $tokenLen = socket_read($this->socket, 2, PHP_BINARY_READ);
        $tokenLen = unpack("n", $tokenLen);
        $tokenLen = $tokenLen[1];
        $token = socket_read($this->socket, $tokenLen, PHP_BINARY_READ);

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