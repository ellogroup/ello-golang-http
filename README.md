# Ello Go HTTP packages

Common packages for handling HTTP requests and returning responses.

## handler

### handler.New

`handler.New()` takes an HTTP handler function with the signature `func(*http.Request) response.Response` and returns a 
function that implements the http.Handler interface for use with the `net/http` package.

This allows for a much simpler handler function that just has to return a response instead of having the responsibility 
of writing that response to the response writer.

## middleware

Some common middleware for use with the `net/http` package.

### middleware.NewAssertJsonPayloadMiddleware

Returns a middleware handler that asserts the HTTP request has a JSON payload by checking the `Content-Type` header. 
Returns a `http.StatusUnsupportedMediaType` (415) response on failure.

### middleware.NewLogCtxMiddleware

Returns a middleware handler that adds [LogCtx](https://github.com/ellogroup/ello-golang-ctx) to the context of the 
request. The `LogCtx` contains context of the request, including method, URI and request ID (if available), which can be 
attached to log entries. After the request has been processed a log entry will be written with additional context 
including status code and response time.

### middleware.NewZapLoggerMiddleware

**Deprecated. Please use LogCtxMiddleware.**

Returns a middleware handler that adds a (zap) logger to the context of the request. The logger contains context of the 
request, including method, URI and request ID (if available), which will be attached to log entries. After the request 
has been processed a log entry will be written with additional context including status code and response time.

### middleware.NewRequestIdMiddleware

Returns a middleware handler that adds a request ID to the request context. This request ID will be from the request 
headers, or generated if not present.

## response

### response.New

Creates a new `Response` from a status code and plain text body.

### response.NewJson

Creates a new JSON `Response` from a status code and an entity to JSON encoded.

### response.NewNoContent

Creates a new `Response` from a status code only.

### response.NewError

Creates a new `ErrorDetails` from a status code. The error code and message are set from the description of the status 
code by default. The error can be configured before transforming into a JSON `Response`.

#### func (ErrorDetails) WithCode(string) ErrorDetails

Set the error code, which should be used to reference the error that has occurred. The error code is _not_ the status 
code.

#### func (ErrorDetails) WithMessage(string) ErrorDetails

Set the error message, which is a human-readable description of the error.

#### func (ErrorDetails) WithMeta(any) ErrorDetails

Set additional metadata. This attribute will be passed through the JSON Marshaller, so structs should include relevant 
`json:` tags.

#### func (ErrorDetails) JsonResponse() Response

Build and return the JSON `Response` from the `ErrorDetails`.

## query

### query.NewValidator

Query validator can validate a map of validation rules against a request query `url.Values`. The validation rules need 
to be supported by `github.com/go-playground/validator`: 
https://pkg.go.dev/github.com/go-playground/validator/v10#readme-fields