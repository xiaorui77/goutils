# HTTPR

http 路由

## 使用

```go
package main

import (
	"github.com/yougtao/goutils/httpr"
	"net/http"
)

func main() {
	router := httpr.NewHttpr()
	router.GET("/", func(c *httpr.Context) {
		c.String(200, "Hello World")
	})

	_ = http.ListenAndServe(":8080", router)
}

```