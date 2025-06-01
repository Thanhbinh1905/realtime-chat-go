package middleware

import (
	"net/http"
	"time"

	"github.com/Thanhbinh1905/realtime-chat/shared/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.LogInfo("Received request", zap.String("method", r.Method), zap.String("url", r.URL.Path))

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logger.LogInfo("Processed request", zap.String("method", r.Method), zap.String("url", r.URL.Path), zap.Duration("duration", duration))
	})
}
