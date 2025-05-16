package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/lulzshadowwalker/personal/internal/config"
	"github.com/lulzshadowwalker/personal/internal/template"
)

func initLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	})
	return slog.New(handler).With("service", "zaya-backend")
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration_ms", duration.Milliseconds(),
			"client_ip", r.RemoteAddr,
		)
	})
}

func main() {
	logger := initLogger()
	port := config.Port()
	logger.Info("loading configuration", "port", port)

	component := template.Hello("zaya")
	logger.Debug("template component ready")

	baseHandler := templ.Handler(component)
	handler := loggingMiddleware(logger, baseHandler)

	addr := fmt.Sprintf(":%s", port)
	logger.Info("starting HTTP server", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Error("server terminated unexpectedly", "error", err)
	}
}
