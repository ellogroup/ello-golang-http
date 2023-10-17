package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

var RequestIDHeader = "X-Request-Id"

type requestIDKey struct{}

var RequestIDCtxKey = &requestIDKey{}

// NewRequestIdMiddleware returns a handler to be used as middleware. This middleware will add a request id to the
// request context. If a request id has been passed in the request headers this will be used, otherwise a random UUID
// will be generated instead.
//
// If used, it is recommended this is one of the first middleware in the chain so all following processes have a
// request id.
func NewRequestIdMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Fetch request id from header
			requestId := r.Header.Get(RequestIDHeader)
			if requestId == "" {
				// Request id not found, generate one
				requestId = uuid.New().String()
			}

			// Add to context
			ctx := context.WithValue(r.Context(), RequestIDCtxKey, requestId)

			// Add to response headers
			w.Header().Set(RequestIDHeader, requestId)

			// Call the next handler in the chain, passing the response writer and
			// the updated request object with the new context value.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestId will extract the request id from the request context. If the request id is not set in the context an empty
// string will be returned.
func RequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if requestID, ok := ctx.Value(RequestIDCtxKey).(string); ok {
		return requestID
	}
	return ""
}
