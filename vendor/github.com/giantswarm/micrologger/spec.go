package micrologger

// Logger is a simple interface describing services that emit messages to gather
// certain runtime information.
type Logger interface {
	// Log takes a sequence of alternating key/value pairs which are used to
	// create the log message structure.
	Log(v ...interface{}) error

	// With returns a new contextual logger with keyvals appended to those
	// passed to calls to Log. If logger is also a contextual logger
	// created by With, keyvals is appended to the existing context.
	With(keyvals ...interface{}) Logger
}
