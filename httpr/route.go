package httpr

import (
	"github.com/xiaorui77/goutils/logx"
	"net/http"
	"strings"
	"time"
)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRouter() *router {
	return &router{
		handlers: map[string]HandlerFunc{},
		roots:    map[string]*node{},
	}
}

// only support simple path routing
func (r *router) registerRoute(method, pattern string, handler HandlerFunc) {
	parts := splitPattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) parseRoute(method, path string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	parts := splitPattern(path)
	no := root.search(parts, 0)
	if no == nil {
		return nil, nil
	}
	ps := splitPattern(no.pattern)
	params := map[string]string{}
	for i, p := range ps {
		if p[0] == ':' {
			params[p[1:]] = parts[i]
		}
		if p[0] == '*' && len(p) > 1 {
			params[p[1:]] = strings.Join(parts[i:], "/")
			break
		}
	}
	return no, params
}

func (r *router) handle(c *Context) {
	no, params := r.parseRoute(c.Method, c.Path)
	if no != nil {
		c.Params = params
		key := c.Method + "-" + no.pattern
		if handler, ok := r.handlers[key]; ok {
			begin := time.Now()
			logx.Infof("[httpr] request [%s] %s - %s", c.RequestId, c.Method, c.Path)
			handler(c)
			logx.Debugf("[httpr] response [%s] complete, cost %s", c.RequestId, time.Now().Sub(begin).String())
		} else {
			logx.Errorf("[httpr] route [%v] parse error", c.Path)
		}
	} else {
		c.StringWithHttpStatus(http.StatusNotFound, "[httpr] 404 NOT FOUND: %s\n", c.Path)
	}
}

func splitPattern(pattern string) []string {
	ps := strings.Split(pattern, "/")

	var parts []string
	for _, p := range ps {
		if p != "" {
			parts = append(parts, p)
			if p[0] == '*' {
				break
			}
		}
	}
	return parts
}
