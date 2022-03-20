package tests

import (
	"bytes"
	"github.com/xiaorui77/goutils/logx"
	"testing"
)

// result: 580 ns/op on Apple M1
func BenchmarkGetCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logx.GetCaller(1)
	}
}

// result: 460 ns/op on Apple M1
func BenchmarkInfo(b *testing.B) {
	var buffer bytes.Buffer
	logx.Init("test", logx.WithOutput(&buffer))
	for i := 0; i < b.N; i++ {
		logx.Infof("Hello %s", "hello")
	}
}

// result: 1800 ns/op on Apple M1
func BenchmarkInfoWithCaller(b *testing.B) {
	var buffer bytes.Buffer
	logger := logx.NewLogx("test")
	logger.SetOutput(&buffer)
	logger.SetReportCaller(true)
	for i := 0; i < b.N; i++ {
		logger.Infof("Hello %s", "hello")
	}
}

// result: 1600 ns/op on Apple M1
func BenchmarkGlobalInfoWithCaller(b *testing.B) {
	var buffer bytes.Buffer
	logx.Init("test", logx.WithOutput(&buffer), logx.WithReportCaller(true))
	for i := 0; i < b.N; i++ {
		logx.Infof("Hello %s", "hello")
	}
}
