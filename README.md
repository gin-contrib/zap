# zap

[![Build Status](https://travis-ci.org/gin-contrib/zap.svg?branch=master)](https://travis-ci.org/gin-contrib/zap) [![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/zap)](https://goreportcard.com/report/github.com/gin-contrib/zap)
[![GoDoc](https://godoc.org/github.com/gin-contrib/zap?status.svg)](https://godoc.org/github.com/gin-contrib/zap)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Alternative logging through [zap](https://github.com/uber-go/zap). Thanks for [Pull Request](https://github.com/gin-gonic/contrib/pull/129) from [@yezooz](https://github.com/yezooz)

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/gin-contrib/zap
```

Import it in your code:

```go
import "github.com/gin-contrib/zap"
```

## Example

See the [example](example/main.go).

[embedmd]:# (example/main.go go)
```go
package main

import (
	"fmt"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.New()

	logger, _ := zap.NewProduction()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	//   - Custom zap filed
	r.Use(ginzap.Logger(logger, ginzap.WithTimeFormat(time.RFC3339),
		ginzap.WithUTC(true),
		ginzap.WithCustomFields(
			func(c *gin.Context) zap.Field { return zap.String("custom_field1", "value1_"+c.ClientIP()) },
			func(c *gin.Context) zap.Field { return zap.String("custom_field2", "value2_"+c.ClientIP()) },
		),
	))
	// simple ginzap.Logger(logger, ginzap.WithTimeFormat(time.RFC3339), ginzap.WithUTC(true))
	// r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	//   - Custom zap filed
	r.Use(ginzap.Recovery(logger, true,
		ginzap.WithCustomFields(
			func(c *gin.Context) zap.Field { return zap.String("custom_field1", "value1_"+c.ClientIP()) },
			func(c *gin.Context) zap.Field { return zap.String("custom_field2", "value2_"+c.ClientIP()) },
		),
	))
	// simple ginzap.Recovery(logger, true)
	// r.Use(ginzap.RecoveryWithZap(logger, true))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example when panic happen.
	r.GET("/panic", func(c *gin.Context) {
		panic("An unexpected error happen!")
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
```
