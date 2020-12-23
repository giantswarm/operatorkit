package sentry

import (
	"context"
	"fmt"
	"reflect"

	"github.com/getsentry/sentry-go"
)

type Default struct {
}

func (s *Default) Capture(ctx context.Context, err error) {
	method := extractReflectedStacktraceMethod(err)
	fmt.Printf("valid: %t, name: %s\n", method.IsValid(), method.String())
	sentry.CaptureException(err)
}

func extractReflectedStacktraceMethod(err error) reflect.Value {
	var method reflect.Value

	// https://github.com/pingcap/errors
	methodGetStackTracer := reflect.ValueOf(err).MethodByName("GetStackTracer")
	// https://github.com/pkg/errors
	methodStackTrace := reflect.ValueOf(err).MethodByName("StackTrace")
	// https://github.com/go-errors/errors
	methodStackFrames := reflect.ValueOf(err).MethodByName("StackFrames")

	if methodGetStackTracer.IsValid() {
		stacktracer := methodGetStackTracer.Call(make([]reflect.Value, 0))[0]
		stacktracerStackTrace := reflect.ValueOf(stacktracer).MethodByName("StackTrace")

		if stacktracerStackTrace.IsValid() {
			method = stacktracerStackTrace
		}
	}

	if methodStackTrace.IsValid() {
		method = methodStackTrace
	}

	if methodStackFrames.IsValid() {
		method = methodStackFrames
	}

	return method
}
