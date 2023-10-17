package mock

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type ResponseWriter struct {
	mock.Mock
}

func (m *ResponseWriter) Header() http.Header {
	args := m.Called()
	h := args.Get(0).(http.Header)
	return h
}

func (m *ResponseWriter) Write(bytes []byte) (int, error) {
	args := m.Called(bytes)
	return args.Int(0), args.Error(1)
}

func (m *ResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}
