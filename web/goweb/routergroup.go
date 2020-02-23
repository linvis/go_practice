package goweb

type RouterGroup struct {
	Base     string
	engine   *Engine
	Handlers HandlerChain
}

func (group *RouterGroup) combineHandlers(handlers HandlerChain) HandlerChain {
	handlers = append(group.Handlers, handlers...)
	return handlers
}

func (group *RouterGroup) combineAbsPath(path string) string {
	if group.Base == "" {
		return path
	}

	if path == "" {
		return group.Base
	}

	if group.Base[len(group.Base)-1] == '/' && path[0] == '/' {
		path = group.Base + path[1:]
	} else {
		path = group.Base + path
	}

	return path
}

func (group *RouterGroup) GET(path string, handlers ...HandlerFunc) {
	path = group.combineAbsPath(path)
	handlers = group.combineHandlers(handlers)

	group.engine.router.addRouter("GET", path, handlers)
}
