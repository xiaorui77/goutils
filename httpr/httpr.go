// Package httpr provide routing and other extended functions.
package httpr

import (
	"net/http"
)

// HandlerFunc defines the request handler used by Context
type HandlerFunc func(c *Context)

// Httpr is core for httpr
type Httpr struct {
	router *router
}

func NewEngine() *Httpr {
	return &Httpr{router: newRouter()}
}

func (e *Httpr) GET(pattern string, handler HandlerFunc) {
	e.addRoute(http.MethodGet, pattern, handler)
}

func (e *Httpr) POST(pattern string, handler HandlerFunc) {
	e.addRoute(http.MethodPost, pattern, handler)
}

func (e *Httpr) PUT(pattern string, handler HandlerFunc) {
	e.addRoute(http.MethodPut, pattern, handler)
}

func (e *Httpr) DELETE(pattern string, handler HandlerFunc) {
	e.addRoute(http.MethodDelete, pattern, handler)
}

func (e *Httpr) addRoute(method, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

func (e *Httpr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.router.handle(c)
}
