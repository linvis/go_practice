package goweb

import "fmt"

type node struct {
	path string
	isWild bool
	isEnd bool
	children []*node
	handlers HandlerChain
}


func initTree() *node {
	return &node{}
}

func longestCommonPrefix(pattern, path string) int {
	i := 0
	for ; i < len(pattern) && i < len(path); i++ {
		if pattern[i] != path[i] {
			break
		}
	}

	return i
}

func (n *node) insertChild(path string, handlers HandlerChain) {
	if len(n.path) == 0 {
		n.path = path
		n.isWild = false
		n.children = []*node{}
		n.handlers = handlers
		n.isEnd = true
	} else {
		child := node{
			path: path,
			isWild: false,
			isEnd: true,
			children: []*node{},
			handlers: handlers,
		}

		n.children = append(n.children, &child)
	}
}

func (n *node) insert(path string, handlers HandlerChain) {
	// fullPath := path

	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, handlers)
		return
	}

	for {
		i := longestCommonPrefix(n.path, path)

		// split 
		if (i < len(n.path)) {
			child := node{
				path : n.path[i:],
				isWild: n.isWild,
				children: n.children,
				handlers: n.handlers,
				isEnd: n.isEnd,
			}

			n.children = []*node{&child}
			n.path = path[:i]
			n.isWild = false
			n.handlers = nil
		}

		if i < len(path) {
			path = path[i:]

			n.insertChild(path, handlers)
			return
		}

		n.handlers = append(n.handlers, handlers...)
		return
	}
}

func (n *node) search(path string) *node {
	if len(n.path) > len(path) {
		return nil
	}

	if n.path != path[:len(n.path)] {
		fmt.Println("no url")
		return nil
	}

	path = path[len(n.path):]
	if len(path) <= 0 && n.isEnd {
		return n
	}

	for _, child := range n.children {
		ans := child.search(path)
		if ans != nil {
			return ans
		}
	}

	return nil
}