package server

import (
	microerror "github.com/giantswarm/microkit/error"
	kithttp "github.com/go-kit/kit/transport/http"
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
		message: errorMessage(config.Underlying),

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

func (e *responseError) IsEndpoint() bool {
	switch u := e.underlying.(type) {
	case kithttp.Error:
		switch u.Domain {
		case kithttp.DomainEncode:
			return true
		case kithttp.DomainDecode:
			return true
		case kithttp.DomainDo:
			return true
		}
	}

	return false
}

func (e *responseError) SetCode(code string) {
	e.code = code
}

func (e *responseError) SetMessage(message string) {
	e.message = message
}

func (e *responseError) Underlying() error {
	kErr, ok := e.underlying.(kithttp.Error)
	if ok {
		return kErr.Err
	}

	return e.underlying
}
