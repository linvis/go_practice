package goweb

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HandlerFunc func(*Context)
type HandlerChain []HandlerFunc

type Engine struct {
	RouterGroup

	htmlTemplates *template.Template // for html render
	router        *router
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
	group := &RouterGroup{
		Base:     path,
		Handlers: HandlerChain{},
		engine:   engine,
	}

	// global middleware
	if len(engine.Handlers) > 0 {
		group.Handlers = append(group.Handlers, engine.Handlers...)
	}

	group.Handlers = append(group.Handlers, handlers...)

	return group
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(template.FuncMap{}).ParseGlob(pattern))
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	method := c.Req.Method
	path := c.Req.URL.Path

	node := engine.router.tree[method].search(path)
	if node == nil {
		c.Status(http.StatusNotFound)
		log.Warn(c.StatusCode, "    ", c.Req.URL.Path)
		return
	}

	c.Method = method
	c.Path = path
	c.Params = node.param
	c.handlers = node.handlers
	c.Next()
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	ctx := &Context{
		Writer: w,
		Req:    req,
		engine: engine,
		index:  -1,
	}

	engine.handleHTTPRequest(ctx)
}

func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
