package micrologger

import "context"

// Logger is a simple interface describing services that emit messages to
// gather certain runtime information.
type Logger interface {
	// Log takes a sequence of alternating key/value pairs which are used
	// to create the log message structure.
	Log(keyVals ...interface{}) error
	// LogCtx is the same as Log but additionally taking a context which
	// may contain additional key-value pairs that are added to the log
	// issuance, if any.
	LogCtx(ctx context.Context, keyVals ...interface{}) error
	// With returns a new contextual logger with keyVals appended to those
	// passed to calls to Log. If logger is also a contextual logger
	// created by With, keyVals is appended to the existing context.
	With(keyVals ...interface{}) Logger
}
