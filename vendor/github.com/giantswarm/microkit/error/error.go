// Package error provides project wide helper functions for a more convenient
// and efficient error handling.
package error

import (
	"fmt"

	"github.com/juju/errgo"
)

var (
	// MaskAny is a simple error masker. Masked errors act as tracers within the
	// source code. Inspecting an masked error shows where the error was passed
	// through within the code base. This is gold for debugging and bug huntin.
	MaskAny = errgo.MaskFunc(errgo.Any)
)

// MaskAnyf is like MaskAny. In addition to that it takes a format string and
// variadic arguments like fmt.Sprintf. The format string and variadic arguments
// are used to annotate the given errgo error.
func MaskAnyf(err error, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	f = fmt.Sprintf("%s: %s", err.Error(), f)
	newErr := errgo.WithCausef(nil, errgo.Cause(err), f, v...)
	newErr.(*errgo.Err).SetLocation(1)

	return newErr
}

// PanicOnError panics in case the given error is not nil. Otherwise it will do
// nothing.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
