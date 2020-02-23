package goweb

import (
	"net/http"
)

type HandlerFunc func(*Context)
type HandlerChain []HandlerFunc

type Engine struct {
	RouterGroup

	ctx    *Context
	router *router
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Base:     "/",
			Handlers: nil,
		},
		router: newRouter(),
	}

	engine.RouterGroup.engine = engine

	return engine
}

func (engine *Engine) Group(path string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Base:     path,
		Handlers: handlers,
		engine:   engine,
	}
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	engine.ctx = newContext(w, req)

	engine.router.handle(engine.ctx)
}

func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
