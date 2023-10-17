package handler

import (
	"context"
	"errors"
	"github.com/ellogroup/ello-golang-http/internal/mock"
	"github.com/ellogroup/ello-golang-http/pkg/http/middleware"
	"github.com/ellogroup/ello-golang-http/pkg/http/response"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name             string
		response         response.Response
		writeBodyError   error
		wantErrorsLogged int
	}{
		{
			name:             "Successfully write response to response writer",
			response:         response.New(http.StatusOK, []byte("test body")),
			wantErrorsLogged: 0,
		},
		{
			name:             "Error logged out when unable to write response",
			response:         response.New(http.StatusOK, []byte("test body")),
			writeBodyError:   errors.New("could not writer to writer error"),
			wantErrorsLogged: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writerMock := new(mock.ResponseWriter)
			headers := http.Header(map[string][]string{})
			writerMock.On("Header").Return(headers)
			writerMock.On("WriteHeader", tt.response.StatusCode).Once()
			writerMock.On("Write", tt.response.BodyEncoded).Return(len(tt.response.BodyEncoded), tt.writeBodyError).Once()

			logCoreMock := new(mock.ZapCore)
			logCoreMock.On("Enabled", zap.ErrorLevel).Return(true).Maybe()

			logMock := zap.New(logCoreMock)

			handler := New(func(*http.Request) response.Response {
				return tt.response
			})

			r := &http.Request{}
			ctx := context.WithValue(r.Context(), middleware.LoggerCtxKey, logMock)
			r = r.WithContext(ctx)

			handler(writerMock, r)

			logCoreMock.AssertNumberOfCalls(t, "Enabled", tt.wantErrorsLogged)
		})
	}
}
