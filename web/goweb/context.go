package goweb

import "net/http"

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Method string
	Path   string
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
	}
}

func (c *Context) String(s string) {
	c.Writer.Write([]byte(s))
}
