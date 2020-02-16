package main

import "fmt"

func main() {
	root := NewNode("", nil, map[string]*Node{})

	routes := []string{
		"/api/users/login",
		"/api/users",
		"/api/user",
		"/api/profiles/:username",
		"/api/profiles/:username/follow",
		"/api/articles",
		"/api/articles/:slug",
		"/api/articles/:slug/comments",
		"/api/articles/:slug/comments/:id",
		"/api/articles/:slug/favorite",
		"/api/articles/feed",
		"/api/tags",
	}
	for _, route := range routes {
		if err := root.insertWithParam(route, nil); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println(root)
}
