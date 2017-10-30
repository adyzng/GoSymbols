package route

import "net/http"

// ResponseLogger is a middleware used to keep an copy of Response.StatusCode.
//
type ResponseLogger struct {
	w          http.ResponseWriter
	StatusCode int
}

// Header returns the header map that will be sent by
// WriteHeader. The Header map also is the mechanism with which
// Handlers can set HTTP trailers.
func (m *ResponseLogger) Header() http.Header {
	return m.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (m *ResponseLogger) Write(data []byte) (int, error) {
	return m.w.Write(data)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (m *ResponseLogger) WriteHeader(status int) {
	if m.StatusCode == 0 {
		// since status code can only write once.
		m.StatusCode = status
	}
	m.w.WriteHeader(status)
}
