package microerror

import (
	"errors"
	"fmt"

	"github.com/juju/errgo"
)

type ErrgoHandlerConfig struct {
	// CallDepth is useful when creating a wrapper for ErrgoHandler. Its
	// value is used to push stack location and skip wrapping function
	// location as an origin. The default value is 0.
	CallDepth int
}

func DefaultErrgoHandlerConfig() ErrgoHandlerConfig {
	return ErrgoHandlerConfig{
		CallDepth: 0,
	}
}

// ErrgoHandler implements Handler interface.
type ErrgoHandler struct {
	callDepth int
	maskFunc  func(err error, allow ...func(error) bool) error
}

func NewErrgoHandler(config ErrgoHandlerConfig) *ErrgoHandler {
	return &ErrgoHandler{
		callDepth: config.CallDepth + 1, // +1 for ErrgoHandler wrapping methods
		maskFunc:  errgo.MaskFunc(errgo.Any),
	}
}

func (h *ErrgoHandler) New(s string) error {
	return errors.New(s)
}

func (h *ErrgoHandler) Newf(f string, v ...interface{}) error {
	return fmt.Errorf(f, v...)
}

func (h *ErrgoHandler) Cause(err error) error {
	return errgo.Cause(err)
}

func (h *ErrgoHandler) Mask(err error) error {
	if err == nil {
		return nil
	}

	newErr := h.maskFunc(err)
	newErr.(*errgo.Err).SetLocation(h.callDepth)
	return newErr
}

func (h *ErrgoHandler) Maskf(err error, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	newErr := errgo.WithCausef(err, errgo.Cause(err), f, v...)
	newErr.(*errgo.Err).SetLocation(h.callDepth)
	return newErr
}
