package goweb

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

func (r *router) addRouter(method string, path string, handlers HandlerChain) {
	// log.Printf("router: %s - %s", method, path)

	r.tree[method].insert(path, handlers)
}
