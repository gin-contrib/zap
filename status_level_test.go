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

func TestStatusLevelMapper(t *testing.T) {
	for _, tc := range []struct {
		name       string
		status     int
		withErrors bool
		mapped     bool
		want       zapcore.Level
	}{
		{"success", http.StatusOK, false, true, zapcore.InfoLevel},
		{"client-error", http.StatusBadRequest, true, true, zapcore.WarnLevel},
		{"server-error", http.StatusInternalServerError, false, true, zapcore.ErrorLevel},
		{"default-errors", http.StatusBadRequest, true, false, zapcore.ErrorLevel},
		{"default-level", http.StatusOK, false, false, zapcore.DebugLevel},
	} {
		t.Run(tc.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.DebugLevel)
			calls := 0
			conf := &Config{DefaultLevel: zapcore.DebugLevel}
			if tc.mapped {
				conf.StatusLevelMapper = func(status int) zapcore.Level {
					calls++
					if status != tc.status {
						t.Errorf("status=%d, want %d", status, tc.status)
					}
					return tc.want
				}
			}
			r := gin.New()
			r.Use(GinzapWithConfig(zap.New(core), conf))
			r.GET("/status", func(c *gin.Context) {
				if tc.withErrors {
					_ = c.Error(errors.New("first"))
					_ = c.Error(errors.New("second"))
				}
				c.Status(tc.status)
			})
			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/status", nil))
			wantCount := 1
			if tc.withErrors {
				wantCount = 2
			}
			if logs.Len() != wantCount {
				t.Fatalf("logs=%d, want %d", logs.Len(), wantCount)
			}
			for _, entry := range logs.All() {
				if entry.Level != tc.want {
					t.Errorf("level=%v, want %v", entry.Level, tc.want)
				}
			}
			if tc.mapped && calls != 1 {
				t.Errorf("mapper called %d times", calls)
			}
		})
	}
}

func TestStatusLevelMapperNotCalledForSkippedPath(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	r := gin.New()
	r.Use(GinzapWithConfig(zap.New(core), &Config{
		SkipPaths:         []string{"/health"},
		StatusLevelMapper: func(int) zapcore.Level { t.Error("mapper called for skipped request"); return zapcore.InfoLevel },
	}))
	r.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/health", nil))
	if logs.Len() != 0 {
		t.Fatal("skipped request was logged")
	}
}
