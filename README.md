# zap

[![Run Tests](https://github.com/gin-contrib/zap/actions/workflows/go.yml/badge.svg)](https://github.com/gin-contrib/zap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/zap)](https://goreportcard.com/report/github.com/gin-contrib/zap)
[![GoDoc](https://godoc.org/github.com/gin-contrib/zap?status.svg)](https://godoc.org/github.com/gin-contrib/zap)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Alternative logging through [zap](https://github.com/uber-go/zap). Thanks for [Pull Request](https://github.com/gin-gonic/contrib/pull/129) from [@yezooz](https://github.com/yezooz)

## Requirement

Require Go **1.19** or later.

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

See the [example](_example/example01/main.go).

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

## Custom Zap fields

example for custom log request body, response request ID or log [Open Telemetry](https://opentelemetry.io/) TraceID.

```go
func main() {
  r := gin.New()

  logger, _ := zap.NewProduction()

  r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
    UTC:        true,
    TimeFormat: time.RFC3339,
    Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
      fields := []zapcore.Field{}
      // log request ID
      if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
        fields = append(fields, zap.String("request_id", requestID))
      }

      // log trace and span ID
      if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
        fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
        fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
      }

      // log request body
      var body []byte
      var buf bytes.Buffer
      tee := io.TeeReader(c.Request.Body, &buf)
      body, _ = io.ReadAll(tee)
      c.Request.Body = io.NopCloser(&buf)
      fields = append(fields, zap.String("body", string(body)))

      return fields
    }),
  }))

  // Example ping request.
  r.GET("/ping", func(c *gin.Context) {
    c.Writer.Header().Add("X-Request-Id", "1234-5678-9012")
    c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  r.POST("/ping", func(c *gin.Context) {
    c.Writer.Header().Add("X-Request-Id", "9012-5678-1234")
    c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    panic(err)
  }
}
```

## Custom `skipper` function

Example for custom `skipper` function

```go
r.Use(GinzapWithConfig(logger, &Config{
  TimeFormat: time.RFC3339,
  UTC: true,
  Skipper: func(c *gin.Context) bool {
    return c.Request.URL.Path == "/ping" && c.Request.Method == "GET"
  },
}))
```

Full example

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

  r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
    UTC:        true,
    TimeFormat: time.RFC3339,
    Skipper: func(c *gin.Context) bool {
      return c.Request.URL.Path == "/ping" && c.Request.Method == "GET"
    },
  }))

  // Example ping request.
  r.GET("/ping", func(c *gin.Context) {
    c.Writer.Header().Add("X-Request-Id", "1234-5678-9012")
    c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  r.POST("/ping", func(c *gin.Context) {
    c.Writer.Header().Add("X-Request-Id", "9012-5678-1234")
    c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    panic(err)
  }
}
```

## Custom `SkipPathRegexps` function

Example for custom `SkipPathRegexps` function

```go
rxURL := regexp.MustCompile(`^/ping\s*`)
r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
  UTC:             true,
  TimeFormat:      time.RFC3339,
  SkipPathRegexps: []*regexp.Regexp{rxURL},
}))
```

Full example

```go
package main

import (
  "fmt"
  "regexp"
  "time"

  ginzap "github.com/gin-contrib/zap"

  "github.com/gin-gonic/gin"
  "go.uber.org/zap"
)

func main() {
  r := gin.New()

  logger, _ := zap.NewProduction()
  rxURL := regexp.MustCompile(`^/ping\s*`)

  r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
    UTC:             true,
    TimeFormat:      time.RFC3339,
    SkipPathRegexps: []*regexp.Regexp{rxURL},
  }))

  // Example ping request.
  r.GET("/ping1234", func(c *gin.Context) {
    c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    panic(err)
  }
}
```
