// Package ginzap provides log handling using zap package.
// Code structure based on ginrus package.
package ginzap

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap on INFO level.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func Ginzap(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return OnLevel(logger, zapcore.InfoLevel, timeFormat, utc)
}

// OnLevel returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap on the provided level.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func OnLevel(logger *zap.Logger, lvl zapcore.Level, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			if ce := logger.Check(lvl, path); ce != nil {
				ce.Write(zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.String("time", end.Format(timeFormat)),
					zap.Duration("latency", latency))
			}
		}
	}
}
