package goweb

import (
	"fmt"
	"testing"
)

func fakeHandler(c *Context) {
	fmt.Println("fakeHandler")
}

func TestTreeInsert(t *testing.T) {

	root := initTree()

	root.insert("/", HandlerChain{fakeHandler})
	root.insert("/he/ha/hb", HandlerChain{fakeHandler})
	root.insert("/he/ha/ff", HandlerChain{fakeHandler})
	root.insert("/he/:name", HandlerChain{fakeHandler})
	root.insert("/static/*filepath", HandlerChain{fakeHandler})

	// ans := root.search("/hello")
	ans := root.search("/he/aaa")
	t.Log(ans)
	ans = root.search("/static/js/info.js")
	t.Log(ans)
}
