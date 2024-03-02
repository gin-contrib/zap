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
