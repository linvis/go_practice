package goweb

import (
	"log"
)

type router struct {
	tree map[string]*node
}

func newRouter() *router {
	r := &router{}
	r.tree = make(map[string]*node)
	r.tree["GET"] = initTree()
	r.tree["POST"] = initTree()
	r.tree["PUT"] = initTree()
	r.tree["DELETE"] = initTree()

	return r
}

func (r *router) addRouter(method string, path string, handler ...HandlerFunc) {
	log.Printf("router: %s - %s", method, path)

	r.tree[method].insert(path, handler)
}

func (r *router) handle(c *Context) {

	n := r.tree[c.Method].search(c.Path)
	if n == nil {
		log.Printf("no path")
		return
	}

	for _, handler := range n.handlers {
		handler(c)
	}
}
