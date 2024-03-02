package ginzap

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const testPath = "/test"

func buildDummyLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, obs := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	return logger, obs
}

func timestampLocationCheck(timestampStr string, location *time.Location) error {
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return err
	}
	if timestamp.Location() != location {
		return fmt.Errorf("timestamp should be utc but %v", timestamp.Location())
	}

	return nil
}

func TestGinzap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := gin.New()

	utcLogger, utcLoggerObserved := buildDummyLogger()
	r.Use(Ginzap(utcLogger, time.RFC3339, true))

	localLogger, localLoggerObserved := buildDummyLogger()
	r.Use(Ginzap(localLogger, time.RFC3339, false))

	r.GET(testPath, func(c *gin.Context) {
		c.JSON(204, nil)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(ctx, "GET", testPath, nil)
	r.ServeHTTP(res1, req1)

	if len(utcLoggerObserved.All()) != 1 {
		t.Fatalf("Log should be 1 line but there're %d", len(utcLoggerObserved.All()))
	}

	logLine := utcLoggerObserved.All()[0]
	pathStr := logLine.Context[2].String
	if pathStr != testPath {
		t.Fatalf("logged path should be /test but %s", pathStr)
	}

	err := timestampLocationCheck(logLine.Context[7].String, time.UTC)
	if err != nil {
		t.Fatal(err)
	}

	if len(localLoggerObserved.All()) != 1 {
		t.Fatalf("Log should be 1 line but there're %d", len(utcLoggerObserved.All()))
	}

	logLine = localLoggerObserved.All()[0]
	pathStr = logLine.Context[2].String
	if pathStr != testPath {
		t.Fatalf("logged path should be /test but %s", pathStr)
	}
}

func TestGinzapWithConfig(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := gin.New()

	utcLogger, utcLoggerObserved := buildDummyLogger()
	r.Use(GinzapWithConfig(utcLogger, &Config{
		TimeFormat:   time.RFC3339,
		UTC:          true,
		SkipPaths:    []string{"/no_log"},
		DefaultLevel: zapcore.WarnLevel,
	}))

	r.GET(testPath, func(c *gin.Context) {
		c.JSON(204, nil)
	})

	r.GET("/no_log", func(c *gin.Context) {
		c.JSON(204, nil)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(ctx, "GET", testPath, nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(ctx, "GET", "/no_log", nil)
	r.ServeHTTP(res2, req2)

	if res2.Code != 204 {
		t.Fatalf("request /no_log is failed (%d)", res2.Code)
	}

	if len(utcLoggerObserved.All()) != 1 {
		t.Fatalf("Log should be 1 line but there're %d", len(utcLoggerObserved.All()))
	}

	logLine := utcLoggerObserved.All()[0]
	pathStr := logLine.Context[2].String
	if pathStr != testPath {
		t.Fatalf("logged path should be /test but %s", pathStr)
	}

	err := timestampLocationCheck(logLine.Context[7].String, time.UTC)
	if err != nil {
		t.Fatal(err)
	}

	if logLine.Level != zapcore.WarnLevel {
		t.Fatalf("log level should be warn but was %s", logLine.Level.String())
	}
}

func TestLoggerSkipper(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := gin.New()

	utcLogger, utcLoggerObserved := buildDummyLogger()
	r.Use(GinzapWithConfig(utcLogger, &Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Skipper: func(c *gin.Context) bool {
			return c.Request.URL.Path == "/no_log"
		},
	}))

	r.GET(testPath, func(c *gin.Context) {
		c.JSON(204, nil)
	})

	r.GET("/no_log", func(c *gin.Context) {
		c.JSON(204, nil)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(ctx, "GET", testPath, nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(ctx, "GET", "/no_log", nil)
	r.ServeHTTP(res2, req2)

	if res2.Code != 204 {
		t.Fatalf("request /no_log is failed (%d)", res2.Code)
	}

	if len(utcLoggerObserved.All()) != 1 {
		t.Fatalf("Log should be 1 line but there're %d", len(utcLoggerObserved.All()))
	}

	logLine := utcLoggerObserved.All()[0]
	pathStr := logLine.Context[2].String
	if pathStr != testPath {
		t.Fatalf("logged path should be /test but %s", pathStr)
	}
}

func TestSkipPathRegexps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := gin.New()

	rxURL := regexp.MustCompile(`^/no_\s*`)

	utcLogger, utcLoggerObserved := buildDummyLogger()
	r.Use(GinzapWithConfig(utcLogger, &Config{
		TimeFormat:      time.RFC3339,
		UTC:             true,
		SkipPathRegexps: []*regexp.Regexp{rxURL},
	}))

	r.GET(testPath, func(c *gin.Context) {
		c.JSON(204, nil)
	})

	r.GET("/no_log", func(c *gin.Context) {
		c.JSON(204, nil)
	})

	res1 := httptest.NewRecorder()
	req1, _ := http.NewRequestWithContext(ctx, "GET", testPath, nil)
	r.ServeHTTP(res1, req1)

	res2 := httptest.NewRecorder()
	req2, _ := http.NewRequestWithContext(ctx, "GET", "/no_log", nil)
	r.ServeHTTP(res2, req2)

	if res2.Code != 204 {
		t.Fatalf("request /no_log is failed (%d)", res2.Code)
	}

	if len(utcLoggerObserved.All()) != 1 {
		t.Fatalf("Log should be 1 line but there're %d", len(utcLoggerObserved.All()))
	}

	logLine := utcLoggerObserved.All()[0]
	pathStr := logLine.Context[2].String
	if pathStr != testPath {
		t.Fatalf("logged path should be /test but %s", pathStr)
	}
}
