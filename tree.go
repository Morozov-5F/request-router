package main

import (
	"errors"
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
	return node.prefix
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
	if l < len(node.prefix) {
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
		if nil == child {
			node.children[str[l:l+1]] = NewNode(str[l:], value, map[string]*Node{})
			return nil
		}
		if err := child.insert(str[l:], value); err != nil {
			return err
		}
	}

	node.value = value
	return nil
}

func (node *Node) getValue(str string) (http.HandlerFunc, error) {
	for len(str) > 0 {
		_, l := commonPrefix(node.prefix, str)
		if l != len(node.prefix) && node.prefix[:1] != ":" {
			return nil, errors.New("no value for given path")
		}

		if l == len(str) {
			return node.value, nil
		}

		child := node.children[str[l:l+1]]
		if nil == child {
			return nil, errors.New("no value for given path")
		}
		return child.getValue(str[l:])
	}
	return nil, errors.New("no value for given path")
}

func splitByParams(str string) []string {
	res := []string{}
	for i := 0; i < len(str); i++ {
		if str[i] == ':' {
			res = append(res, str[:i])
			for i < len(str) && str[i] != '/' {
				i++
			}
			res = append(res, str[:i])
		}
	}
	if len(res) == 0 {
		res = append(res, str)
	}
	return res
}

func (node *Node) insertWithParam(str string, value http.HandlerFunc) error {
	params := splitByParams(str)
	for i, path := range params {
		var newValue http.HandlerFunc = nil
		if i == len(params)-1 {
			newValue = value
		}
		if err := node.insert(path, newValue); err != nil {
			return err
		}
	}
	return nil
}
