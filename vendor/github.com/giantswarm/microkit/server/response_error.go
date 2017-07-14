package server

import (
	microerror "github.com/giantswarm/microkit/error"
)

// ResponseErrorConfig represents the configuration used to create a new
// response error.
type ResponseErrorConfig struct {
	// Settings.
	Underlying error
}

// DefaultResponseErrorConfig provides a default configuration to create a new
// response error by best effort.
func DefaultResponseErrorConfig() ResponseErrorConfig {
	return ResponseErrorConfig{
		// Settings.
		Underlying: nil,
	}
}

// New creates a new configured response error.
func NewResponseError(config ResponseErrorConfig) (ResponseError, error) {
	// Settings.
	if config.Underlying == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "underlying must not be empty")
	}

	newResponseError := &responseError{
		// Internals.
		code:    CodeInternalError,
		message: config.Underlying.Error(),

		// Settings.
		underlying: config.Underlying,
	}

	return newResponseError, nil
}

type responseError struct {
	// Internals.
	code    string
	message string

	// Settings.
	underlying error
}

func (e *responseError) Code() string {
	return e.code
}

func (e *responseError) Error() string {
	return e.underlying.Error()
}

func (e *responseError) Message() string {
	return e.message
}

func (e *responseError) SetCode(code string) {
	e.code = code
}

func (e *responseError) SetMessage(message string) {
	e.message = message
}

func (e *responseError) Underlying() error {
	return e.underlying
}
