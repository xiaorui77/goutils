package logx

import "testing"

// result: 660 ns/op
func BenchmarkWithGetCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getCaller()
	}
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Infof("Hello %s", "hello")
	}
}
