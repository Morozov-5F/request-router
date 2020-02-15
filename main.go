package main

import "fmt"

func main() {
	root := NewNode("", nil, map[string]*Node{})
	root.insert("/api/users/login", nil)
	root.insert("/api/users", nil)
	root.insert("/api/user", nil)
	root.insert("/api/profiles/:username", nil)
	root.insert("/api/profiles/:username/follow", nil)
	root.insert("/api/articles", nil)
	root.insert("/api/articles/feed", nil)
	root.insert("/api/articles/:slug", nil)
	root.insert("/api/articles/:slug/comments", nil)
	root.insert("/api/articles/:slug/comments/:id", nil)
	root.insert("/api/articles/:slug/favorite", nil)
	root.insert("/api/tags", nil)

	fmt.Println(root)
}
