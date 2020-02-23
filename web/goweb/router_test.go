package goweb

import (
	"testing"
)

func TestAddRouter(t *testing.T) {
	r := newRouter()
	r.addRouter("GET", "/hello", nil)
	r.addRouter("GET", "/he", nil)
	r.addRouter("GET", "/he/:name", nil)
	r.addRouter("GET", "/he/*file", nil)

	t.Log(r)
}
