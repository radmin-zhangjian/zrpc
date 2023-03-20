package zrpc

import "log"

type methodTree struct {
	method string
	root   *node
}

type methodTrees []methodTree

type node struct {
	path     string
	children []*node
	handlers handlersChain
	fullPath string
}

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

func (n *node) addRoute(pattern string, handlers handlersChain) {
	if pattern == "/" {
		n.path = "/"
		n.handlers = handlers
		return
	}

	for k, v := range n.children {
		if v.path == pattern {
			h := v.handlers
			n.children[k].handlers = append(handlers, h[len(h)-1:]...)
			log.Println("route repeat value")
			return
		}
	}

	child := node{
		fullPath: pattern,
		path:     pattern,
		handlers: handlers,
	}
	n.children = append(n.children, &child)
}

func (n *node) getRoute(pattern string) (handlers handlersChain) {
	if pattern == "/" {
		handlers = n.handlers
		return
	}
	for k, v := range n.children {
		if v.path == pattern {
			handlers = n.children[k].handlers
			return
		}
	}
	return
}
