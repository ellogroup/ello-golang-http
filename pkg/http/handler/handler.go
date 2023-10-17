package handler

import (
	"github.com/ellogroup/ello-golang-http/pkg/http/middleware"
	"github.com/ellogroup/ello-golang-http/pkg/http/response"
	"go.uber.org/zap"
	"net/http"
)

// New converts a function that takes a request and returns a response, and returns a new handler function that
// implements the http.Handler interface for a http server. This wrapper will pass the request to the handler, and then
// write the response to the response writer.
func New(handler func(r *http.Request) response.Response) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := handler(r)
		if err := resp.WriteTo(w); err != nil {
			// Unable to write the response to the response writer
			log := middleware.Logger(r.Context())
			log.Error("Unable to write response", zap.Error(err))
		}
	}
}
