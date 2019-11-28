package microerror

import (
	"fmt"

	"github.com/juju/errgo"
)

// Stack prints the error with the stack if its argument is underlying
// microerror error or result of Error function otherwise. Its main purpose is
// to be used for a value for "stack" micrologger key.
//
// Example:
//
//	logger.LogCtx(ctx, "level", "error", "message", "failed to do a thing", "stack", microerror.Stack(err))
//
func Stack(err error) string {
	switch err.(type) {
	case nil:
		return fmt.Sprintf("%v", nil)
	case *errgo.Err:
		return fmt.Sprintf("%#v", err)
	default:
		return err.Error()
	}
}

// Desc returns the description of a microerror.Error.
func Desc(err error) string {
	c := Cause(err)
	switch c.(type) {
	case nil:
		return ""
	case *Error:
		e, ok := c.(*Error)
		if ok {
			return e.Desc
		}
		return ""
	default:
		return ""
	}
}
