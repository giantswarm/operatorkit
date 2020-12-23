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
	pcs := extractPcs(method)
	fmt.Println(pcs)
	sentry.CaptureException(err)
}

func extractPcs(method reflect.Value) []uintptr {
	var pcs []uintptr

	stacktrace := method.Call(make([]reflect.Value, 0))[0]

	if stacktrace.Kind() != reflect.Slice {
		return nil
	}

	for i := 0; i < stacktrace.Len(); i++ {
		pc := stacktrace.Index(i)

		if pc.Kind() == reflect.Uintptr {
			pcs = append(pcs, uintptr(pc.Uint()))
			continue
		}

		if pc.Kind() == reflect.Struct {
			field := pc.FieldByName("ProgramCounter")
			if field.IsValid() && field.Kind() == reflect.Uintptr {
				pcs = append(pcs, uintptr(field.Uint()))
				continue
			}
		}
	}

	return pcs
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
		fmt.Println("methodGetStackTracer")
		stacktracer := methodGetStackTracer.Call(make([]reflect.Value, 0))[0]
		stacktracerStackTrace := reflect.ValueOf(stacktracer).MethodByName("StackTrace")

		if stacktracerStackTrace.IsValid() {
			method = stacktracerStackTrace
		}
	}

	if methodStackTrace.IsValid() {
		fmt.Println("methodStackTrace")
		method = methodStackTrace
	}

	if methodStackFrames.IsValid() {
		fmt.Println("methodStackFrames")
		method = methodStackFrames
	}

	return method
}
