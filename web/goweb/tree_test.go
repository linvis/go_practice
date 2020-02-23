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
	root.insert("/he/:name/bb", HandlerChain{fakeHandler})
	root.insert("/static/*filepath", HandlerChain{fakeHandler})
	root.insert("/user/:name", HandlerChain{fakeHandler})

	// ans := root.search("/hello")
	ans := root.search("/user/gordon")
	ans = root.search("/he/aa/bb")
	t.Log(ans)
	ans = root.search("/static/js/info.js")
	t.Log(ans)
}
