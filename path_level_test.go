package ginzap

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestPathLevels(t *testing.T) {
	for _, tc := range []struct {
		name, target string
		want         zapcore.Level
		fail, skip   bool
	}{
		{"override", "/health", zapcore.DebugLevel, false, false},
		{"query", "/health?probe=1", zapcore.DebugLevel, false, false},
		{"exact-path", "/health/other", zapcore.InfoLevel, false, false},
		{"unmatched", "/other", zapcore.InfoLevel, false, false},
		{"error", "/health", zapcore.ErrorLevel, true, false},
		{"skip", "/health", zapcore.DebugLevel, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			levels := map[string]zapcore.Level{"/health": zapcore.DebugLevel}
			conf := &Config{PathLevels: levels}
			if tc.skip {
				conf.SkipPaths = []string{"/health"}
			}
			r := gin.New()
			r.Use(GinzapWithConfig(zap.New(core), conf))
			// Configuration ownership follows SkipPaths: later caller mutation is ignored.
			levels["/health"] = zapcore.ErrorLevel
			r.NoRoute(func(c *gin.Context) {
				c.Request.URL.Path = "/rewritten"
				if tc.fail {
					_ = c.Error(errors.New("failed"))
				}
				c.Status(http.StatusOK)
			})
			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, tc.target, nil))
			if tc.skip {
				if logs.Len() != 0 {
					t.Fatal("skipped path logged")
				}
				return
			}
			if logs.Len() != 1 {
				t.Fatalf("logs=%d", logs.Len())
			}
			if got := logs.All()[0].Level; got != tc.want {
				t.Errorf("level=%v, want %v", got, tc.want)
			}
		})
	}
}
