package main

import (
	"fmt"
)

func main() {

	tree := NewTree()
	tree.Insert("/api/users/login", 1)
	tree.Insert("/api/users", 2)
	tree.Insert("/api/profiles/:username", 3)
	tree.Insert("/api/profiles/:username/follow", 4)
	tree.Insert("/api/articles", 5)
	tree.Insert("/api/articles/feed", 6)
	tree.Insert("/api/articles/:slug", 7)
	tree.Insert("/api/articles/:slug/comments", 8)
	tree.Insert("/api/articles/:slug/comments/:id", 9)
	tree.Insert("/api/articles/:slug/favorite", 10)
	tree.Insert("/api/tags", 11)

	fmt.Printf("%v\n", tree)
}
