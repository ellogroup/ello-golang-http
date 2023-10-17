package response

import (
	"net/http"
	"strings"
)

type ErrorBody struct {
	ErrorDetails ErrorDetails `json:"error"`
}
type ErrorDetails struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Meta    any    `json:"meta"`
}

func (e ErrorDetails) WithCode(c string) ErrorDetails {
	return ErrorDetails{
		Status:  e.Status,
		Code:    c,
		Message: e.Message,
		Meta:    e.Meta,
	}
}

func (e ErrorDetails) WithMessage(m string) ErrorDetails {
	return ErrorDetails{
		Status:  e.Status,
		Code:    e.Code,
		Message: m,
		Meta:    e.Meta,
	}
}

func (e ErrorDetails) WithMeta(m any) ErrorDetails {
	return ErrorDetails{
		Status:  e.Status,
		Code:    e.Code,
		Message: e.Message,
		Meta:    m,
	}
}

func (e ErrorDetails) JsonResponse() Response {
	return NewJson(e.Status, ErrorBody{
		ErrorDetails: e,
	})
}

// NewError creates a new JSON response from a status code.
func NewError(statusCode int) ErrorDetails {
	return ErrorDetails{
		Status:  statusCode,
		Code:    slugify(http.StatusText(statusCode)),
		Message: http.StatusText(statusCode),
		Meta:    nil,
	}
}

func slugify(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
