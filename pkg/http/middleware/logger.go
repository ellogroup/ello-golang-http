package middleware

import (
	"context"
	"errors"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type loggerKey struct{}

var LoggerCtxKey = &loggerKey{}

// NewZapLoggerMiddleware returns a handler to be used as middleware. This middleware will add details of the request to
// the logger, and add the logger to the context of the request. This logger can then be extracted from the request
// context and any log entries will include the details of the request. Once the request is complete, the details of the
// completed request will also be logged out.
//
// If used, it is recommended this is one of the first middleware in the chain so all following processes have access
// to a logger with the request details set. However, the request id middleware should always come _before_ this
// middleware.
func NewZapLoggerMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add context to logger
			requestLog := log.With(
				zap.String("http_proto", r.Proto),
				zap.String("http_method", r.Method),
				zap.String("request_uri", r.RequestURI),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			// Add request id to logger context
			if requestId := RequestId(r.Context()); requestId != "" {
				requestLog = requestLog.With(zap.String("request_id", requestId))
			}

			// Log request info
			requestLog.Info("Request started")

			// Add logger to context
			ctx := context.WithValue(r.Context(), LoggerCtxKey, requestLog)

			// Wrap the response writer, so we can access details of the response, such as status code
			ww := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Log the response details
			requestStart := time.Now()
			defer func() {
				// Log response info
				requestLog.Info(
					"Request complete",
					zap.Int("status_code", ww.Status()),
					zap.Int("duration_ms", int(time.Since(requestStart).Milliseconds())),
					zap.Int("bytes_written", ww.BytesWritten()),
				)
			}()

			// Call the next handler in the chain, passing the response writer and
			// the updated request object with the new context value.
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

// LoggerOrError will extract the logger from the request context. If the logger is not set in the context an error
// will be returned.
func LoggerOrError(ctx context.Context) (*zap.Logger, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if log, ok := ctx.Value(LoggerCtxKey).(*zap.Logger); ok {
		return log, nil
	}
	return nil, errors.New("logger not found within context")
}

// Logger will extract the logger from the request context. If the logger is not set in the context a noop logger will
// be returned.
func Logger(ctx context.Context) *zap.Logger {
	if log, err := LoggerOrError(ctx); err == nil {
		return log
	}
	return zap.NewNop()
}
