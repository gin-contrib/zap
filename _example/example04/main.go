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
