package middleware

import (
	"github.com/ellogroup/ello-golang-ctx/logctx"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// TODO: Comment description
// TODO: Tests
// NewLogCtxMiddleware TODO
func NewLogCtxMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add attributes to context
			ctx := logctx.Add(
				r.Context(),
				logctx.String("http_proto", r.Proto),
				logctx.String("http_method", r.Method),
				logctx.String("request_uri", r.RequestURI),
				logctx.String("remote_addr", r.RemoteAddr),
				logctx.String("user_agent", r.UserAgent()),
			)

			// Add request id to logger context
			if requestId := RequestId(r.Context()); requestId != "" {
				ctx = logctx.Add(ctx, logctx.String("request_id", requestId))
			}

			// Log request info
			log.Info("Request started", logctx.Zap(ctx)...)

			// Wrap the response writer, so we can access details of the response, such as status code
			ww := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Log the response details
			requestStart := time.Now()
			defer func() {
				// Log response info
				log.Info(
					"Request complete",
					logctx.Zap(
						ctx,
						zap.Int("status_code", ww.Status()),
						zap.Int("bytes_written", ww.BytesWritten()),
						zap.Duration("duration", time.Since(requestStart)),
					)...,
				)
			}()

			// Call the next handler in the chain, passing the response writer and
			// the updated request object with the new context value.
			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
