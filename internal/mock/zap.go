package mock

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

type ZapCore struct {
	mock.Mock
}

func (m *ZapCore) Enabled(level zapcore.Level) bool {
	args := m.Called(level)
	return args.Bool(0)
}

func (m *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	args := m.Called(fields)
	c := args.Get(0).(zapcore.Core)
	return c
}

func (m *ZapCore) Check(_ zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce
}

func (m *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	args := m.Called(entry, fields)
	return args.Error(0)
}

func (m *ZapCore) Sync() error {
	args := m.Called()
	return args.Error(0)
}
