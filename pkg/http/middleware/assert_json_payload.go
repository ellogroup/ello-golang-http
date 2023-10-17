package middleware

import (
	"github.com/ellogroup/ello-golang-http/pkg/http/response"
	"go.uber.org/zap"
	"net/http"
)

// NewAssertJsonPayloadMiddleware returns a handler to be used as middleware. This middleware will assert the request
// content type is set to json, otherwise it will return an error and prevent the request from being processed.
func NewAssertJsonPayloadMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != "application/json" {
				log := Logger(r.Context())
				log.Debug("Unexpected Content-Type provided", zap.String("Content-Type", r.Header.Get("Content-Type")))
				if err := response.NewError(http.StatusUnsupportedMediaType).JsonResponse().WriteTo(w); err != nil {
					// Unable to write the response to the response writer
					log.Error("Unable to write response", zap.Error(err))
				}
				return
			}

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}
