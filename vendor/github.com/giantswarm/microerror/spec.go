package microerror

type Handler interface {
	// New returns a new error with the given error message. It is
	// a drop-in replacement for errors.New from the standard library.
	//
	// NOTE deprecated
	//
	New(s string) error

	// Newf returns a new error with the given printf-formatted error
	// message.
	//
	// NOTE deprecated
	//
	Newf(f string, v ...interface{}) error

	// Cause returns the cause of the given error. If the cause of the err can not
	// be found it returns the err itself.
	//
	// Cause is the usual way to diagnose errors that may have been wrapped by Mask
	// or Maskf.
	Cause(err error) error

	// Mask is a simple error masker. Masked errors act as tracers within the
	// source code. Inspecting an masked error shows where the error was passed
	// through within the code base. This is gold for debugging and bug hunting.
	Mask(err error) error

	// Maskf is like Mask. In addition to that it takes a format string and
	// variadic arguments like fmt.Sprintf. The format string and variadic
	// arguments are used to annotate the given error.
	Maskf(err error, f string, v ...interface{}) error
}
