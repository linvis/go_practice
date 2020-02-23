package goweb

import (
	"net/http"
	"path"
)

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

func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.Handlers = append(group.Handlers, handlers...)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absPath := path.Join(group.Base, relativePath)
	fileServer := http.StripPrefix(absPath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	group.GET(urlPattern, handler)
}
