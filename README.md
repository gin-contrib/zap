# zap

[![Run Tests](https://github.com/gin-contrib/zap/actions/workflows/go.yml/badge.svg)](https://github.com/gin-contrib/zap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/zap)](https://goreportcard.com/report/github.com/gin-contrib/zap)
[![GoDoc](https://godoc.org/github.com/gin-contrib/zap?status.svg)](https://godoc.org/github.com/gin-contrib/zap)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Alternative logging through [zap](https://github.com/uber-go/zap). Thanks for [Pull Request](https://github.com/gin-gonic/contrib/pull/129) from [@yezooz](https://github.com/yezooz)

## Usage

### Start using it

Download and install it:

```sh
go get github.com/gin-contrib/zap
```

Import it in your code:

```go
import "github.com/gin-contrib/zap"
```

## Example

See the [example](example/main.go).

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
  r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

  // Logs all panic to error log
  //   - stack means whether output the stack info.
  r.Use(ginzap.RecoveryWithZap(logger, true))

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

## Skip logging
When you want to skip logging for specific path,
please use GinzapWithConfig

```go

r.Use(GinzapWithConfig(utcLogger, &Config{
    TimeFormat: time.RFC3339,
    UTC: true,
    SkipPaths: []string{"/no_log"},
}))
```
