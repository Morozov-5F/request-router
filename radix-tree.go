package main

import (
	"fmt"
	"net/http"
)

type node struct {
	children map[string]*node
	prefix   string
	value    int

	handler http.HandlerFunc
}

func (n *node) String() string {
	return n.dump("", "")
}

func (n *node) Prefix() string {
	if n.prefix != "" {
		return fmt.Sprintf("%v (%v)", n.prefix, n.value)
	}
	return "root"
}

func (n *node) dump(prefix string, childPrefix string) string {
	res := prefix + n.Prefix() + "\n"
	var keys []string
	for k := range n.children {
		keys = append(keys, k)
	}
	for i, k := range keys {
		if i == len(keys)-1 {
			res += n.children[k].dump(childPrefix+"└── ", childPrefix+"    ")
		} else {
			res += n.children[k].dump(childPrefix+"├── ", childPrefix+"|   ")
		}
	}
	return res
}

func newNode(prefix string, value int) *node {
	return &node{prefix: prefix, children: make(map[string]*node), value: value}
}

type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{root: newNode("", 0)}
}

func (n *Tree) String() string {
	return n.root.String()
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func getCommonPrefix(a string, b string) string {
	for i := 0; i < len(a); i++ {
		if i == len(b) || a[i] != b[i] {
			return a[:i]
		}
	}
	return a
}

func (n *node) insert(prefix string, value int) {
	if prefix == n.prefix {
		n.value = value
		return
	}

	commonPrefix := getCommonPrefix(prefix, n.prefix)
	if commonPrefix == n.prefix {
		toAdd := prefix[len(commonPrefix):]
		lookup := n.children[toAdd[:1]]
		if lookup == nil {
			n.children[toAdd[:1]] = newNode(toAdd, value)
		} else {
			lookup.insert(toAdd, value)
		}

		return
	}

	newPrefix := n.prefix[len(commonPrefix):]
	nodeToAdd := newNode(newPrefix, n.value)
	nodeToAdd.children = n.children

	n.children = map[string]*node{newPrefix[:1]: nodeToAdd}
	n.prefix = commonPrefix
	n.value = 0

	if commonPrefix == prefix {
		n.value = value
		return
	}

	toAdd := prefix[len(commonPrefix):]
	n.children[toAdd[:1]] = newNode(toAdd, value)
}

func (n *Tree) Insert(prefix string, value int) {
	n.root.insert(prefix, value)
}
