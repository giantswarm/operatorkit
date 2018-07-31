// Package microerror provides project wide helper functions for a more convenient
// and efficient error handling.
package microerror

var (
	handler Handler = NewErrgoHandler(ErrgoHandlerConfig{
		CallDepth: 1,
	})
)

// New returns a new error with the given error message. It is a drop-in
// replacement for errors.New from the standard library.
func New(s string) error {
	return handler.New(s)
}

// Newf returns a new error with the given printf-formatted error message.
func Newf(f string, v ...interface{}) error {
	return handler.Newf(f, v...)
}

// Cause returns the cause of the given error. If the cause of the err can not
// be found it returns the err itself.
//
// Cause is the usual way to diagnose errors that may have been wrapped by Mask
// or Maskf.
func Cause(err error) error {
	return handler.Cause(err)
}

// Mask is a simple error masker. Masked errors act as tracers within the
// source code. Inspecting an masked error shows where the error was passed
// through within the code base. This is gold for debugging and bug hunting.
func Mask(err error) error {
	return handler.Mask(err)
}

// Maskf is like Mask. In addition to that it takes a format string and
// variadic arguments like fmt.Sprintf. The format string and variadic
// arguments are used to annotate the given errgo error.
func Maskf(err error, f string, v ...interface{}) error {
	return handler.Maskf(err, f, v...)
}
