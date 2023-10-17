package response

import (
	"encoding/json"
	"net/http"
)

const contentTypePlainText = "text/plain"
const contentTypeJson = "application/json"

type Response struct {
	StatusCode  int
	BodyEncoded []byte
	BodyDecoded any
	ContentType string
	Headers     http.Header
}

// New creates a new plain text response.
func New(statusCode int, body []byte) Response {
	return Response{
		StatusCode:  statusCode,
		BodyEncoded: body,
		ContentType: contentTypePlainText,
		Headers:     map[string][]string{},
	}
}

// NewJson creates a new JSON response.
func NewJson(statusCode int, body any) Response {
	return Response{
		StatusCode:  statusCode,
		BodyDecoded: body,
		ContentType: contentTypeJson,
		Headers:     map[string][]string{},
	}
}

// NewNoContent creates a new response with no body.
func NewNoContent(statusCode int) Response {
	return Response{
		StatusCode: statusCode,
		Headers:    map[string][]string{},
	}
}

func (r Response) Clone() Response {
	return Response{
		StatusCode:  r.StatusCode,
		BodyEncoded: r.BodyEncoded,
		BodyDecoded: r.BodyDecoded,
		ContentType: r.ContentType,
		Headers:     r.Headers.Clone(),
	}
}

func (r Response) WithContentType(ct string) Response {
	c := r.Clone()
	c.ContentType = ct
	return c
}

func (r Response) WithHeader(key, value string) Response {
	c := r.Clone()
	c.Headers.Add(key, value)
	return c
}

func (r Response) WriteTo(w http.ResponseWriter) error {
	// Write headers
	if r.Headers != nil {
		for k, v := range r.Headers {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	}
	if len(r.ContentType) > 0 && len(r.Headers.Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", r.ContentType)
	}

	// Write status code
	w.WriteHeader(r.StatusCode)

	// Write body
	if r.ContentType == contentTypeJson {
		// JSON
		if r.BodyDecoded != nil {
			return json.NewEncoder(w).Encode(r.BodyDecoded)
		}
	}

	if len(r.BodyEncoded) > 0 {
		// Plain
		_, err := w.Write(r.BodyEncoded)
		return err
	}

	return nil
}
