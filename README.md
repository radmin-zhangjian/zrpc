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
