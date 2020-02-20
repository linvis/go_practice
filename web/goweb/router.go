package goweb

import (
	"net/http"

	"log"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRouter(method string, url string, handler HandlerFunc) {
	key := method + url
	log.Printf("router: %s - %s", method, url)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + c.Path

	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusOK, "404 not found")
	}
}
