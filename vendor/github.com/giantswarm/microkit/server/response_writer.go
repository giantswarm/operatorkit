package server

import (
	"bytes"
	"net/http"

	microerror "github.com/giantswarm/microkit/error"
)

// ResponseWriterConfig represents the configuration used to create a new
// response writer.
type ResponseWriterConfig struct {
	// Settings.
	BodyBuffer     *bytes.Buffer
	ResponseWriter http.ResponseWriter
	StatusCode     int
}

// DefaultResponseWriterConfig provides a default configuration to create a new
// response writer by best effort.
func DefaultResponseWriterConfig() ResponseWriterConfig {
	return ResponseWriterConfig{
		// Settings.
		BodyBuffer:     &bytes.Buffer{},
		ResponseWriter: nil,
		StatusCode:     http.StatusOK,
	}
}

// New creates a new configured response writer.
func NewResponseWriter(config ResponseWriterConfig) (ResponseWriter, error) {
	// Settings.
	if config.BodyBuffer == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "body buffer must not be empty")
	}
	if config.ResponseWriter == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "response writer must not be empty")
	}
	if config.StatusCode == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "status code must not be empty")
	}

	newResponseWriter := &responseWriter{
		// Settings.
		bodyBuffer:     config.BodyBuffer,
		responseWriter: config.ResponseWriter,
		statusCode:     config.StatusCode,
	}

	return newResponseWriter, nil
}

type responseWriter struct {
	// Settings.
	bodyBuffer     *bytes.Buffer
	responseWriter http.ResponseWriter
	statusCode     int
}

func (rw *responseWriter) BodyBuffer() *bytes.Buffer {
	return rw.bodyBuffer
}

func (rw *responseWriter) Header() http.Header {
	return rw.responseWriter.Header()
}

func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	_, err := rw.bodyBuffer.Write(b)
	if err != nil {
		return 0, microerror.MaskAny(err)
	}

	return rw.responseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(c int) {
	rw.responseWriter.WriteHeader(c)
	rw.statusCode = c
}
