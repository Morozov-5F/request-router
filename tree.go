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

func (node *Node) insert(str string, value http.HandlerFunc) {
	p, l := commonPrefix(node.prefix, str)

	// If the p length is less than existing p, then we should split the node
	if l < len(node.prefix) {
		node.children = map[string]*Node{node.prefix[l : l+1]: NewNode(node.prefix[l:], node.value, node.children)}
		node.prefix = p
		node.value = nil
	}

	if l < len(str) {
		// Find a child that starts with the same symbol as new value
		child := node.children[str[l:l+1]]
		if nil == child {
			node.children[str[l:l+1]] = NewNode(str[l:], value, map[string]*Node{})
			return
		}
		child.insert(str[l:], value)
		return
	}

	node.value = value
}

func (node *Node) getValue(str string) (http.HandlerFunc, error) {
	for len(str) > 0 {
		_, l := commonPrefix(node.prefix, str)
		if l != len(node.prefix) {
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
