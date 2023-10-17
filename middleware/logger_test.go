package middleware

import (
	"context"
	"fmt"
	"github.com/ellogroup/ello-golang-http/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestLoggerOrError(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	logCoreMock := new(mock.ZapCore)
	logger := zap.New(logCoreMock)
	tests := []struct {
		name    string
		args    args
		want    *zap.Logger
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Logger within context returns successfully",
			args:    args{context.WithValue(context.Background(), LoggerCtxKey, logger)},
			want:    logger,
			wantErr: assert.NoError,
		},
		{
			name:    "Non-logger within context returns error",
			args:    args{context.WithValue(context.Background(), LoggerCtxKey, "not a logger")},
			wantErr: assert.Error,
		},
		{
			name:    "No logger within context returns error",
			args:    args{context.Background()},
			wantErr: assert.Error,
		},
		{
			name:    "nil context returns error",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoggerOrError(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("LoggerOrError(%v)", tt.args.ctx)) {
				return
			}
			assert.Equalf(t, tt.want, got, "LoggerOrError(%v)", tt.args.ctx)
		})
	}
}

func TestLogger(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	logCoreMock := new(mock.ZapCore)
	logger := zap.New(logCoreMock)
	tests := []struct {
		name string
		args args
		want *zap.Logger
	}{
		{
			name: "Logger in context returns successfully",
			args: args{context.WithValue(context.Background(), LoggerCtxKey, logger)},
			want: logger,
		},
		{
			name: "Logger not in context returns noop logger",
			args: args{context.Background()},
			want: zap.NewNop(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Logger(tt.args.ctx), "Logger(%v)", tt.args.ctx)
		})
	}
}
