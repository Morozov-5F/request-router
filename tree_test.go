package main

import "testing"

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
