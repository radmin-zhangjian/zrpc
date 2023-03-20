package zrpc

import (
	"context"
	"log"
	"math"
	"reflect"
	"strings"
	"sync"
)

const abortIndex int8 = math.MaxInt8 / 2

type handlerFunc func(*Context)
type handlersChain []handlerFunc

type Context struct {
	Ctx        context.Context
	Args       any
	Reply      *any
	inArgs     []reflect.Value
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

func serverMothed(pattern string) (serviceName string, methodName string) {
	dot := strings.LastIndex(pattern, ".")
	if dot < 0 {
		log.Fatalf("rpc: service/method request ill-formed: %s", pattern)
		return
	}
	serviceName = pattern[:dot]
	methodName = pattern[dot+1:]
	return
}

func (group *RouterGroup) addRoute(pattern string, handlers handlersChain) {
	//assert1(pattern[0] == '/', "path must begin with '/'")
	assert1(len(handlers) > 0, "there must be at least one handler")

	group.mu.Lock()
	defer group.mu.Unlock()

	log.Printf("full path: %s", pattern)

	serviceName, path := serverMothed(pattern)
	root := group.trees.get(serviceName)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		group.trees = append(group.trees, methodTree{method: serviceName, root: root})
	}
	root.addRoute(pattern, path, handlers)
}

func (group *RouterGroup) GetRoute(pattern string) handlersChain {
	serviceName, path := serverMothed(pattern)
	root := group.trees.get(serviceName)
	if root == nil {
		log.Printf("rpc: service/method not find: %s", pattern)
	}
	return root.getRoute(pattern, path)
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
