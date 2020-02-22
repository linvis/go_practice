package goweb

import (
	"fmt"
	"strings"
)

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

func findWildcard(path string) (string, int, bool) {
	valid := false

	for i, c := range path {
		if c != ':' && c != '*' {
			continue
		}

		valid = true
		for end := i+1; end < len(path); end++ {
			if path[end] == '/' {
				return path[i:end], i, valid
			}
			if path[end] == '*' || path[end] == ':' {
				return "", i, false
			}
		}

		return path[i:], i, valid
	}

	return "", 0, false
}

func (n *node) insertChild(path string, handlers HandlerChain) {

	if (len(n.path)) != 0 {
		child := &node{}
		n.children = append(n.children, child)
		n = child
	}

	wild, idx, valid := findWildcard(path)
	if !valid {
		n.path = path
		n.isWild = false
		n.children = []*node{}
		n.handlers = handlers
		n.isEnd = true
		return
	}

	n.path = path[:idx]
	n.children = []*node{}

	child := &node{
		path : wild,
		children : []*node{},
		isWild: true,
	}
	n.children = append(n.children, child)
	n = child

	child = &node{
		path: path[idx+len(wild):],
		children : []*node{},
		isEnd: true,
		handlers: handlers,
	}
	n.children = append(n.children, child)
	n = child

	return
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

type nodeInfo struct {
	handlers HandlerChain
	param map[string]string
}

func (n *node) search(path string) *nodeInfo {

	param := strings.Split(path, "/")[0]
	paramMap := make(map[string]string, 0)

	if n.isWild && len(param) > 0 {
		path = path[len(param):]
		paramMap[n.path[1:]] = param
		goto LOOP
	}

	if len(n.path) > len(path) {
		return nil
	}

	if n.path != path[:len(n.path)] {
		fmt.Println("no url")
		return nil
	}

	path = path[len(n.path):]
	if len(path) <= 0 && n.isEnd {
		return &nodeInfo{
			handlers: n.handlers,
			param: paramMap,
		}
	}

LOOP:
	for _, child := range n.children {
		info := child.search(path)
		if info != nil {
			for k, v := range paramMap {
				info.param[k] = v
			}

			return info
		}
	}

	return nil
}