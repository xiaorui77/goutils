package demo

import (
	"github.com/xiaorui77/goutils/httpr"
	"net/http"
	"time"
)

// httpr use demo.
func _() {
	router := httpr.NewEngine()
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
