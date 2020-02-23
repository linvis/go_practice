package goweb

import (
	"strings"
)

type node struct {
	path     string
	isWild   bool
	isEnd    bool
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
		for end := i + 1; end < len(path); end++ {
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

	if idx > 0 {
		n.path = path[:idx]
		n.children = []*node{}

		child := &node{
			path:     wild,
			children: []*node{},
			isWild:   true,
		}
		n.children = append(n.children, child)
		n = child
	} else {
		n.path = wild
		n.children = []*node{}
		n.isWild = true
	}

	if idx+len(wild) < len(path) {
		child := &node{
			path:     path[idx+len(wild):],
			children: []*node{},
			isEnd:    true,
			handlers: handlers,
		}
		n.children = append(n.children, child)
		n = child
	} else {
		n.isEnd = true
		n.handlers = handlers
	}

	return
}

func (n *node) insert(path string, handlers HandlerChain) {
	// fullPath := path

	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, handlers)
		return
	}

LOOP:
	for {
		i := longestCommonPrefix(n.path, path)

		// split
		// update current node's path to prefix
		// insert reset as a child
		if i < len(n.path) {
			child := node{
				path:     n.path[i:],
				isWild:   n.isWild,
				children: n.children,
				handlers: n.handlers,
				isEnd:    n.isEnd,
			}

			n.children = []*node{&child}
			n.path = path[:i]
			n.isWild = false
			n.handlers = nil
		}

		// insert new path's reset
		if i < len(path) {
			path = path[i:]

			for _, child := range n.children {
				if len(child.path) > 0 && child.path[0] == path[0] {
					n = child
					goto LOOP
				}
			}

			n.insertChild(path, handlers)
			return
		}

		n.handlers = append(n.handlers, handlers...)
		n.isEnd = true
		return
	}
}

type nodeInfo struct {
	handlers HandlerChain
	param    map[string]string
}

func (n *node) search(path string) *nodeInfo {

	param := strings.Split(path, "/")[0]
	paramMap := make(map[string]string, 0)
	pattern := n.path

	if !n.isWild {
		if len(pattern) > len(path) {
			return nil
		}

		// not equal
		if pattern != path[:len(pattern)] {
			return nil
		}

		// find node
		path = path[len(pattern):]
		if len(path) <= 0 && n.isEnd {
			return &nodeInfo{
				handlers: n.handlers,
				param:    paramMap,
			}
		}
	} else {
		// empty param
		if len(param) <= 0 {

		}

		if pattern[0] == ':' {
			paramMap[pattern[1:]] = param
			path = path[len(param):]

			if len(path) <= 0 {
				return &nodeInfo{
					handlers: n.handlers,
					param:    paramMap,
				}
			}
		} else if pattern[0] == '*' {
			paramMap[pattern[1:]] = path

			if n.isEnd {
				return &nodeInfo{
					handlers: n.handlers,
					param:    paramMap,
				}
			} else {
				return nil
			}
		}
	}

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
