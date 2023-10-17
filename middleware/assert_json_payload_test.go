package middleware

import (
	"context"
	"errors"
	"github.com/ellogroup/ello-golang-http/internal/mock"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestNewAssertJsonPayloadMiddleware(t *testing.T) {
	tests := []struct {
		name                  string
		header                http.Header
		writeBodyError        error
		wantStatusCodeWritten int
		wantNextCalled        bool
		wantLogCount          int
	}{
		{
			name:                  "No headers writes StatusUnsupportedMediaType",
			header:                map[string][]string{},
			wantStatusCodeWritten: http.StatusUnsupportedMediaType,
			wantNextCalled:        false,
			wantLogCount:          1,
		},
		{
			name: "Invalid Content Type writes StatusUnsupportedMediaType",
			header: map[string][]string{
				"Content-Type": {"text/plain"},
			},
			wantStatusCodeWritten: http.StatusUnsupportedMediaType,
			wantNextCalled:        false,
			wantLogCount:          1,
		},
		{
			name: "Invalid Content Type and unable to write status logs error",
			header: map[string][]string{
				"Content-Type": {"text/plain"},
			},
			writeBodyError:        errors.New("could not writer to writer error"),
			wantStatusCodeWritten: http.StatusUnsupportedMediaType,
			wantNextCalled:        false,
			wantLogCount:          1,
		},
		{
			name: "Valid Content Type calls next in chain",
			header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantNextCalled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewAssertJsonPayloadMiddleware()

			nextCalled := false
			next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				nextCalled = true
			})

			writerMock := new(mock.ResponseWriter)
			headers := http.Header(map[string][]string{})
			writerMock.On("Header").Return(headers)
			writerMock.On("WriteHeader", tt.wantStatusCodeWritten).Once()
			writerMock.On("Write", testifymock.Anything).Return(1, nil).Maybe()

			logCoreMock := new(mock.ZapCore)
			logCoreMock.On("Enabled", zap.DebugLevel).Return(true).Maybe()
			logCoreMock.On("Enabled", zap.ErrorLevel).Return(true).Maybe()

			logMock := zap.New(logCoreMock)

			r := &http.Request{
				Header: tt.header,
			}
			ctx := context.WithValue(r.Context(), LoggerCtxKey, logMock)
			r = r.WithContext(ctx)

			sut(next).ServeHTTP(writerMock, r)

			assert.Equalf(t, tt.wantNextCalled, nextCalled, "nextCalled")

			logCoreMock.AssertNumberOfCalls(t, "Enabled", tt.wantLogCount)
		})
	}
}
