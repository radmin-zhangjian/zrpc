package codec

// Response 定义RPC交互的数据结构
type Response struct {
	// 访问的函数
	ServiceMethod string
	// 访问时的参数
	Args any
	// 返回数据
	Reply any
	// 错误
	Error error
	// Id
	Seq uint64
}
