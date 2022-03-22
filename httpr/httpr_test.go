package httpr

import (
	"github.com/xiaorui77/goutils/logx"
	"net/http"
	"testing"
	"time"
)

func startServer() {
	router := NewEngine()
	router.GET("/hello", func(c *Context) {
		c.String("Hello World!")
	})
	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       15 * time.Second,
	}
	_ = server.ListenAndServe()
}

func TestHttpr_GET(t *testing.T) {
	logx.SetLevel(logx.DebugLevel)
	go startServer()
	time.Sleep(time.Second)
	_, _ = http.Get("http://127.0.0.1:8080/hello")
	time.Sleep(time.Second * 3)
}
