package main

import (
	"net/http"
	"reflect"
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
	routes := map[string]http.HandlerFunc{
		"/api/users/login":                 handlerGenerator("/api/users/login"),
		"/api/users":                       handlerGenerator("/api/users"),
		"/api/user":                        handlerGenerator("/api/user"),
		"/api/articles":                    handlerGenerator("/api/articles"),
		"/api/articles/:slug":              handlerGenerator("/api/articles/:slug"),
		"/api/articles/:slug/comments":     handlerGenerator("/api/articles/:slug/comments"),
		"/api/articles/:slug/comments/:id": handlerGenerator("/api/articles/:slug/comments/:id"),
		"/api/tags":                        handlerGenerator("/api/tags"),
		"/test":                            handlerGenerator("/test"),
	}

	root := NewNode("", nil, map[string]*Node{})
	for route, handler := range routes {
		if err := root.insert(route, handler); err != nil {
			t.Errorf("insertion failed: %v", err)
		}
	}
	type args struct {
		str string
	}
	type want struct {
		cb     http.HandlerFunc
		cbOut  string
		params map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{"Get handler for the first-level route succeeds", args{"/test"}, want{routes["/test"], "/test", nil}, false},
		{"Get handler for the second-level route succeeds", args{"/api/user"}, want{routes["/api/user"], "/api/user", nil}, false},
		{"Get handler for the third-level route succeeds", args{"/api/articles/feed"}, want{routes["/api/articles/:slug"], "/api/articles/:slug", map[string]string{"slug": "feed"}}, false},
		{"Get handler for the fifth-level route succeeds", args{"/api/articles/feed/comments/123456"}, want{routes["/api/articles/:slug/comments/:id"], "/api/articles/:slug/comments/:id", map[string]string{"slug": "feed", "id": "123456"}}, false},
		{"Get handler for non-existing route fails", args{"/non-existing"}, want{nil, "", nil}, true},
		{"Get handler for existing route part fails", args{"/api"}, want{nil, "", nil}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, params, err := root.getValue(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("error is expected but handler is not nil (got = %v)", got)
				return
			}
			if got == nil && tt.want.cb == nil && params == nil && tt.want.params == nil {
				return
			}
			if !reflect.DeepEqual(tt.want.params, params) {
				t.Errorf("wrong parameters for route '%v': expected '%v', got '%v' from channel", tt.args.str, tt.want.params, params)
			}
			if err == nil {
				got(nil, nil)
			}
			select {
			case prefix := <-handlerChannel:
				if prefix != tt.want.cbOut {
					t.Errorf("wrong callback for route '%v', expected '%v', got '%v' from channel", tt.args.str, tt.want.cbOut, prefix)
				}
			case <-time.After(1 * time.Second):
				t.Errorf("timeout while waiting for the callback value for route '%v'", tt.args.str)
			}
		})
	}
}
