// Package loggermeta stores and accesses the container struct in
// context.Context.
package loggermeta

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// loggerMeta is the key for logger struct values in context.Context. Clients use
// loggermeta.NewContext and loggermeta.FromContext instead of using this
// key directly.
var loggerMetaKey key = "loggerMeta"

// LoggerMeta is a communication structure used to transport information in order
// for a micro logger to use it when issuing logs.
type LoggerMeta struct {
	// KeyVals is a mapping of key-value pairs a micro logger adds to the log
	// message issuance.
	KeyVals map[string]string
}

func New() *LoggerMeta {
	return &LoggerMeta{
		KeyVals: map[string]string{},
	}
}

// NewContext returns a new context.Context that carries value v.
func NewContext(ctx context.Context, v *LoggerMeta) context.Context {
	if v == nil {
		return ctx
	}

	return context.WithValue(ctx, loggerMetaKey, v)
}

// FromContext returns the logger struct, if any.
func FromContext(ctx context.Context) (*LoggerMeta, bool) {
	v, ok := ctx.Value(loggerMetaKey).(*LoggerMeta)
	return v, ok
}
