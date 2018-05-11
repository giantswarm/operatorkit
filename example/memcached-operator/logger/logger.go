package logger

import (
	"context"
	"fmt"

	"github.com/giantswarm/micrologger"
)

var Default micrologger.Logger

func init() {
	var err error

	Default, err = micrologger.New(micrologger.Config{})
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func Log(keyVals ...interface{}) error {
	return Default.Log(keyVals...)
}

func LogCtx(ctx context.Context, keyVals ...interface{}) error {
	return Default.LogCtx(ctx, keyVals...)
}
