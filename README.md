# zap

[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/zap)](https://goreportcard.com/report/github.com/gin-contrib/zap)
[![GoDoc](https://godoc.org/github.com/gin-contrib/zap?status.svg)](https://godoc.org/github.com/gin-contrib/zap)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Alternative logging through [zap](https://github.com/uber-go/zap). Thanks for [Pull Request](https://github.com/gin-gonic/contrib/pull/129) from [@yezooz](https://github.com/yezooz)
## Example

See the [example](example/main.go).

```go
package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/zap"
	"github.com/uber-go/zap"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	log := zap.New(
		zap.NewTextEncoder(),
		zap.DebugLevel,
	)

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

```
