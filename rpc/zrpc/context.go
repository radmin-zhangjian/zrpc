package zrpc

import (
	"context"
	"log"
	"math"
	"reflect"
	"sync"
)

const abortIndex int8 = math.MaxInt8 / 2

type handlerFunc func(*Context)
type handlersChain []handlerFunc

type Context struct {
	Rcvr   reflect.Value
	Ctx    context.Context
	Args   any
	Reply  *any
	inArgs []reflect.Value

	StatusCode int
	Keys       map[string]any
	index      int8
	handlers   handlersChain
}

type RouterGroup struct {
	handlersMap map[string]handlersChain
	//handlersMap sync.Map
	handlers handlersChain
	mu       sync.Mutex
	trees    methodTrees
}

func NewContext() *Context {
	return &Context{}
}

func NewRouterGroup() *RouterGroup {
	return &RouterGroup{}
}

func (c *Context) reset() {
	c.handlers = nil
	c.index = -1
	c.StatusCode = 0
	c.Keys = nil
	c.Rcvr = reflect.Value{}
	c.Ctx = nil
	c.Args = nil
	c.Reply = nil
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) Test(h handlersChain) {
	c.reset()
	c.handlers = h
	c.Next()
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func (group *RouterGroup) addRoute(pattern string, handlers handlersChain) {
	//assert1(pattern[0] == '/', "path must begin with '/'")
	assert1(len(handlers) > 0, "there must be at least one handler")

	group.mu.Lock()
	defer group.mu.Unlock()

	log.Printf("Route %s", pattern)
	//group.handlersMap[keys] = handlers
	//keys := strings.Split(pattern, "/")
	//keys = keys[1:]
	//for k, v := range keys {
	//	log.Printf("keys key %d, value %s", k, v)
	//	if k == len(keys)-1 {
	//		group.handlersMap[v] = handlers
	//		//if _, dup := group.handlersMap.LoadOrStore(v, handlers); dup {
	//		//	log.Fatalf("rpc: service already defined: %s", v)
	//		//}
	//	} else {
	//		group.handlersMap[v] = make(map[string]any)
	//		//if _, dup := group.handlersMap.LoadOrStore(v, make(map[string]any)); dup {
	//		//	log.Fatalf("rpc: service already defined: %s", v)
	//		//}
	//	}
	//}
	//log.Printf("group.handlersMap %+v", group.handlersMap)
	////val, ok := group.handlersMap.Load("")
	//val, ok := group.handlersMap["v1"]
	//if ok {
	//	log.Printf("group.handlersMap %+v", val)
	//}

	root := group.trees.get("post")
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		group.trees = append(group.trees, methodTree{method: "post", root: root})
	}
	root.addRoute(pattern, handlers)
}

func (group *RouterGroup) GetRoute(pattern string) handlersChain {
	root := group.trees.get("post")
	if root == nil {

	}
	return root.getRoute(pattern)
}

func (group *RouterGroup) Use(middleware ...handlerFunc) {
	group.handlers = append(group.handlers, middleware...)
}

func (group *RouterGroup) UseHandle(pattern string, handler ...handlerFunc) {
	finalSize := len(group.handlers) + len(handler)
	mergedHandlers := make(handlersChain, finalSize)
	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handler)
	//mergedHandlers = append(mergedHandlers, handler...)
	group.addRoute(pattern, mergedHandlers)
}
