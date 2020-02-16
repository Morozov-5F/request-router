package main

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
)

type Node struct {
	children map[string]*Node

	prefix string
	value  http.HandlerFunc
}

func NewNode(prefix string, value http.HandlerFunc, children map[string]*Node) *Node {
	return &Node{prefix: prefix, value: value, children: children}
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func commonPrefix(a string, b string) (string, int) {
	limit := min(len(a), len(b))
	for i := 0; i < limit; i++ {
		if a[i] != b[i] {
			return a[:i], i
		}
	}
	return a[:limit], limit
}

func (node *Node) repr() string {
	if len(node.prefix) == 0 {
		return "root"
	}
	return fmt.Sprintf("%v (%v)", node.prefix, node.value)
}

func (node *Node) dump(pfx string, childPfx string) string {
	res := pfx + " " + node.repr() + "\n"

	var keys []string
	for k := range node.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i == len(keys)-1 {
			res += node.children[k].dump(childPfx+"└─── ", childPfx+"      ")
			continue
		}
		res += node.children[k].dump(childPfx+"├─── ", childPfx+"│     ")
	}

	return res
}

func (node *Node) String() string {
	return node.dump("", "")
}

func (node *Node) insert(str string, value http.HandlerFunc) error {
	p, l := commonPrefix(node.prefix, str)

	// If the p length is less than existing p, then we should split the node
	if l < len(node.prefix) && l != 0 {
		node.children = map[string]*Node{node.prefix[l : l+1]: NewNode(node.prefix[l:], node.value, node.children)}
		node.prefix = p
		node.value = nil
	}

	if l < len(str) {
		// Find a child that starts with the same symbol as new value
		firstChar := str[l : l+1]
		if nil != node.children[":"] && firstChar != ":" {
			return errors.New("unable to register regular route -- parametric route is already present")
		}
		if nil == node.children[":"] && firstChar == ":" && len(node.children) != 0 {
			return errors.New("unable to register parametric route -- regular route is already present")
		}

		child := node.children[firstChar]
		if nil != child {
			return child.insert(str[l:], value)
		}

		for i := l; i < len(str); i++ {
			if i != l && str[i] == ':' {
				node.children[firstChar] = NewNode(str[l:i], nil, map[string]*Node{})
				return node.children[firstChar].insert(str[i:], value)
			}
			if str[l] == ':' && str[i] == '/' {
				node.children[firstChar] = NewNode(str[l:i], nil, map[string]*Node{})
				return node.children[firstChar].insert(str[i:], value)
			}
		}
		node.children[firstChar] = NewNode(str[l:], value, map[string]*Node{})
		return nil
	}

	node.value = value
	return nil
}

func (node *Node) getValue(str string) (handler http.HandlerFunc, params map[string]string, err error) {
	for len(str) > 0 {
		_, l := commonPrefix(node.prefix, str)
		if l != len(node.prefix) && node.prefix[:1] != ":" {
			return nil, nil, errors.New("no value for given path")
		}

		if len(node.prefix) > 0 && node.prefix[:1] == ":" {
			i := 0
			for ; i < len(str); i++ {
				if str[i] == '/' {
					break
				}
			}

			if i == len(str) && node.value == nil {
				return nil, nil, errors.New("no value for given path")
			}

			if i != len(str) && len(node.children) == 0 {
				return nil, nil, errors.New("no value for given path")
			}

			l = i
			params = map[string]string{node.prefix[1:]: str[0:l]}
		}

		if l == len(str) {
			return node.value, params, nil
		}

		child := node.children[str[l:l+1]]
		if nil == child {
			child = node.children[":"]
			if nil == child {
				return nil, nil, errors.New("no value for given path")
			}
		}
		h, p, e := child.getValue(str[l:])
		if e != nil {
			return nil, nil, e
		}
		if params == nil {
			return h, p, e
		}
		for k, v := range p {
			params[k] = v
		}
		return h, params, nil
	}
	return nil, nil, errors.New("no value for given path")
}
