package logger

// Logger is a simple interface describing services that emit messages to gather
// certain runtime information.
type Logger interface {
	// Log takes a sequence of alternating key/value pairs which are used to
	// create the log message structure.
	Log(v ...interface{}) error
}
