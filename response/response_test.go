package response

import (
	"errors"
	"github.com/ellogroup/ello-golang-http/internal/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type decodedPayload struct {
	Example string `json:"example"`
}

func TestResponse_WriteTo(t *testing.T) {
	type fields struct {
		statusCode  int
		bodyEncoded []byte
		bodyDecoded any
		contentType string
		headers     http.Header
	}
	tests := []struct {
		name           string
		fields         fields
		writeBodyError error
		wantHeader     http.Header
		wantBody       []byte
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Headers are written correctly",
			fields: fields{
				statusCode: 200,
				headers: map[string][]string{
					"header-one": {"One"},
					"header-two": {"Two", "Three"},
				},
			},
			wantHeader: map[string][]string{
				"Header-One": {"One"},
				"Header-Two": {"Two", "Three"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Content-Type in headers is not overwritten by response attribute",
			fields: fields{
				statusCode:  200,
				contentType: "content-type-response",
				headers: map[string][]string{
					"header-one":   {"One"},
					"Content-Type": {"content-type-header"},
					"header-two":   {"Two", "Three"},
				},
			},
			wantHeader: map[string][]string{
				"Header-One":   {"One"},
				"Content-Type": {"content-type-header"},
				"Header-Two":   {"Two", "Three"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Content-Type missing from headers is replaced by response attribute",
			fields: fields{
				statusCode:  200,
				contentType: "content-type-response",
				headers: map[string][]string{
					"header-one": {"One"},
					"header-two": {"Two", "Three"},
				},
			},
			wantHeader: map[string][]string{
				"Header-One":   {"One"},
				"Content-Type": {"content-type-response"},
				"Header-Two":   {"Two", "Three"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Decoded payload JSON-encoded when content type is json",
			fields: fields{
				statusCode:  200,
				contentType: "application/json",
				bodyDecoded: decodedPayload{Example: "test"},
			},
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: []byte("{\"example\":\"test\"}\n"),
			wantErr:  assert.NoError,
		},
		{
			name: "Decoded payload not JSON-encoded when content type is not json",
			fields: fields{
				statusCode:  200,
				contentType: "not/json",
				bodyDecoded: decodedPayload{Example: "test"},
			},
			wantHeader: map[string][]string{
				"Content-Type": {"not/json"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Empty decoded payload not JSON-encoded when content type is json",
			fields: fields{
				statusCode:  200,
				contentType: "application/json",
			},
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Write error when JSON-encoding returns error",
			fields: fields{
				statusCode:  200,
				contentType: "application/json",
				bodyDecoded: decodedPayload{Example: "test"},
			},
			writeBodyError: errors.New("write error"),
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: []byte("{\"example\":\"test\"}\n"),
			wantErr:  assert.Error,
		},
		{
			name: "Encoded payload written when content type is json and no decoded payload set",
			fields: fields{
				statusCode:  200,
				contentType: "application/json",
				bodyEncoded: []byte("test-123"),
			},
			wantHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantBody: []byte(`test-123`),
			wantErr:  assert.NoError,
		},
		{
			name: "Encoded payload written when content type is not json",
			fields: fields{
				statusCode:  200,
				bodyEncoded: []byte("test-123"),
			},
			wantHeader: map[string][]string{},
			wantBody:   []byte(`test-123`),
			wantErr:    assert.NoError,
		},
		{
			name: "Write error when encoded payload written returns error",
			fields: fields{
				statusCode:  200,
				bodyEncoded: []byte("test-123"),
			},
			writeBodyError: errors.New("write error"),
			wantHeader:     map[string][]string{},
			wantBody:       []byte("test-123"),
			wantErr:        assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Response{
				StatusCode:  tt.fields.statusCode,
				BodyEncoded: tt.fields.bodyEncoded,
				BodyDecoded: tt.fields.bodyDecoded,
				ContentType: tt.fields.contentType,
				Headers:     tt.fields.headers,
			}

			writerMock := new(mock.ResponseWriter)
			headers := http.Header(map[string][]string{})
			writerMock.On("Header").Return(headers).Maybe()
			writerMock.On("WriteHeader", tt.fields.statusCode).Once()
			writerMock.On("Write", tt.wantBody).Return(1, tt.writeBodyError).Maybe()

			if !tt.wantErr(t, r.WriteTo(writerMock), "WriteTo()") {
				return
			}

			assert.Equalf(t, tt.wantHeader, headers, "WriteTo().Header")
		})
	}
}
