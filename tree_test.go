package main

import (
	"net/http"
	"testing"
	"time"
)

func Test_commonPrefix(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
	}{
		{"Prefix correct (b is longer)", args{"test", "tester"}, "test", 4},
		{"Prefix correct (a is longer)", args{"tester", "test"}, "test", 4},
		{"Prefix correct (a is not contained in b)", args{"tea", "test"}, "te", 2},
		{"Prefix correct (b is not contained in a)", args{"test", "tea"}, "te", 2},
		{"Prefix correct (a is empty)", args{"", "tester"}, "", 0},
		{"Prefix correct (b is empty)", args{"tester", ""}, "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := commonPrefix(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("commonPrefix() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("commonPrefix() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTree_getValue(t *testing.T) {
	// It's impossible to compare the function address in Go so I borrowed an idea from https://stackoverflow.com/a/51885524
	// every handler should spit the unique value to this channel from its handler
	handlerChannel := make(chan string, 1)
	handlerGenerator := func(retVal string) http.HandlerFunc {
		return func(http.ResponseWriter, *http.Request) {
			handlerChannel <- retVal
		}
	}
	regularRoutes := map[string]http.HandlerFunc{
		"/api/users/login":   handlerGenerator("/api/users/login"),
		"/api/users":         handlerGenerator("/api/users"),
		"/api/user":          handlerGenerator("/api/user"),
		"/api/articles":      handlerGenerator("/api/articles"),
		"/api/articles/feed": handlerGenerator("/api/articles/feed"),
		"/api/tags":          handlerGenerator("/api/tags"),
		"/test":              handlerGenerator("/test"),
	}

	root := NewNode("", nil, map[string]*Node{})
	for route, handler := range regularRoutes {
		root.insert(route, handler)
	}
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    http.HandlerFunc
		wantErr bool
	}{
		{"Get handler for the first-level route succeeds", args{"/test"}, regularRoutes["/test"], false},
		{"Get handler for the second-level route succeeds", args{"/api/user"}, regularRoutes["/api/user"], false},
		{"Get handler for the third-level route succeeds", args{"/api/articles/feed"}, regularRoutes["/api/articles/feed"], false},
		{"Get handler for non-existing route fails", args{"/non-existing"}, nil, true},
		{"Get handler for existing route part fails", args{"/api"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := root.getValue(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("error is expected but handler is not nil (got = %v)", err, got)
				return
			}
			if got == nil && tt.want == nil {
				return
			}
			got(nil, nil)
			select {
			case prefix := <-handlerChannel:
				if prefix != tt.args.str {
					t.Errorf("wrong callback for route '%v', expected '%v', got '%v' from channel", tt.args.str, prefix, tt.args.str)
				}
			case <-time.After(1 * time.Second):
				t.Errorf("timeout while waiting for the callback value for route '%v'", tt.args.str)
			}
		})
	}
}
