package goweb

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	ctx    *Context
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.router.addRouter("GET", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	engine.ctx = newContext(w, req)

	engine.router.handle(engine.ctx)
}

func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
