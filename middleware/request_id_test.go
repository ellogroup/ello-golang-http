package middleware

import (
	"context"
	"github.com/ellogroup/ello-golang-http/internal/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

func TestNewRequestIdMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		header    http.Header
		wantMatch *regexp.Regexp
	}{
		{
			name: "Request id passed in header is used",
			header: map[string][]string{
				RequestIDHeader: {"test-abc"},
			},
			wantMatch: regexp.MustCompile("^test-abc$"),
		},
		{
			name:      "Request id not passed in header a UUID is generated",
			header:    map[string][]string{},
			wantMatch: regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewRequestIdMiddleware()

			var ctx context.Context
			next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				ctx = r.Context()
			})

			headers := http.Header(map[string][]string{})
			writerMock := new(mock.ResponseWriter)
			writerMock.On("Header").Return(headers).Once()

			r := &http.Request{
				Header: tt.header,
			}

			sut(next).ServeHTTP(writerMock, r)

			if !assert.Regexpf(t, tt.wantMatch, headers.Get(RequestIDHeader), "NewRequestIdMiddleware() response header") {
				return
			}

			assert.Regexpf(t, tt.wantMatch, ctx.Value(RequestIDCtxKey).(string), "NewRequestIdMiddleware() context")
		})
	}
}

func TestRequestId(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Request ID within context returns successfully",
			args: args{ctx: context.WithValue(context.Background(), RequestIDCtxKey, "test-123")},
			want: "test-123",
		},
		{
			name: "Non-string within context returns empty string",
			args: args{ctx: context.WithValue(context.Background(), RequestIDCtxKey, 123)},
			want: "",
		},
		{
			name: "No request ID within context returns empty string",
			args: args{ctx: context.Background()},
			want: "",
		},
		{
			name: "Nil context returns empty string",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RequestId(tt.args.ctx); got != tt.want {
				t.Errorf("RequestId() = '%v', want '%v'", got, tt.want)
			}
			assert.Equalf(t, tt.want, RequestId(tt.args.ctx), "RequestId(%v)", tt.args.ctx)
		})
	}
}
