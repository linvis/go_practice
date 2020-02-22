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

	root.insert("/hello", HandlerChain{fakeHandler})
	root.insert("/he", HandlerChain{fakeHandler})

	ans := root.search("/hello")
	t.Log(ans)
}