# grpc
***
### 安装grpc
````
go get google.golang.org/grpc  
````

### 安装protoc
````
brew install protoc
````

### 安装针对go的protoc插件
方法1，使用go install <module>@latest安装  
````
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  
````

方法2，在一个已经包含go.mod文件的项目里使用go get <module>  
````
go get google.golang.org/protobuf/cmd/protoc-gen-go  
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc  
````

安装成功后，会在$GOPATH/bin目录下生成两个2进制文件  
````
protoc-gen-go*  
protoc-gen-go-grpc*  
````

#### 配置环境变量  
````
# 复制文件到 /usr/local/bin/
cp protoc-gen-go /usr/local/bin/
cp protoc-gen-go-grpc /usr/local/bin/

# 创建 bash_profile
vim ~/.bash_profile  
export GOPATH=$HOME/go PATH=$PATH:$GOPATH/bin  
source ~/.bash_profile  
````

### 编译 protoc  
````
protoc --go_out=. --go-grpc_out=. ./hello.proto
protoc -I . --go_out=. --go-grpc_out=. ./response.proto
# 或者
protoc --go_out=. ./hello.proto
protoc --go-grpc_out=. ./hello.proto
# .pb.go文件路径不依赖于.proto文件中的option go_package配置项，直接在go_out指定的目录下生成.pb.go文件
protoc --go_opt=paths=source_relative --go_out=. ./hello.proto
protoc --go-grpc_opt=paths=source_relative --go-grpc_out=. ./hello.proto
````
