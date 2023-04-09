### zrpc
带有中间件的服务

### server 服务端
````
// 注册服务
func register(srv *zrpc.Server) {
	// 将服务端方法，注册一下
	srv.RegisterName(new(service.Test), "service")
	srv.RegisterName(new(v1.Test), "v1")
	srv.RegisterName(new(v2.Test), "v2")
	// 中间件 实例  Use放在最前面时，所有注册的服务端方法都会生效
	srv.Use(func(c *zrpc.Context) {
		log.Println("srv use ==========")
		c.Next()
		log.Println("srv use next ==========")
	})
	srv.UseHandle("v1.QueryInt", func(c *zrpc.Context) {
		log.Println("v1.QueryInt ++++++++++")
		log.Println("c.Args ++++++++++", c.Args)
	}, func(c *zrpc.Context) {
		log.Println("v1.QueryInt ----------")
		c.Next()
		log.Println("c.QueryInt next ----------")
	})
	srv.UseHandle("v1.QueryIntC", func(c *zrpc.Context) {
		log.Println("v1.QueryIntC ++++++++++")
		c.Next()
		log.Println("v1.QueryIntC c.Args ++++++++++", c.Args)
	})
	srv.UseHandle("service.QueryUserC", func(c *zrpc.Context) {
		log.Println("service.QueryUserC ++++++++++")
		c.Next()
		log.Println("service.QueryUserC c.Args ++++++++++", c.Args)
	})

	// 创建中间件
	//r := rpc.NewRouterGroup()
	//r.Use(func(c *rpc.Context) {
	//	log.Println("/Use")
	//	c.Next()
	//	log.Println("/Use next")
	//})
	//r.UseHandle("/v1/test", func(c *rpc.Context) {
	//	log.Println("/v1/test111")
	//	c.Next()
	//	log.Println("/v1/test111-222")
	//}, func(c *rpc.Context) {
	//	log.Println("/v1/test222")
	//})
	//r.UseHandle("/v2/test", func(c *rpc.Context) {
	//	log.Println("/v2/test")
	//})
	//h := r.GetRoute("/v1/test")
	//log.Printf("handles: %+v\n", h)
	//context := rpc.NewContext()
	//context.Test(h)
}

````