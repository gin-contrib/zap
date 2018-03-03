package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.New()

	cfg := zap.NewProductionConfig()
	logger, _ := cfg.Build()
	currentLevel := cfg.Level

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Or use it to log only when on at least DebugLevel
	r.Use(ginzap.OnLevel(logger, zap.DebugLevel, time.RFC3339, true))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	r.GET("/level", gin.WrapF(currentLevel.ServeHTTP))
	r.PUT("/level", gin.WrapF(currentLevel.ServeHTTP))

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
