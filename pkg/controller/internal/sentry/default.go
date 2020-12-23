package sentry

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getsentry/sentry-go"
)

type Default struct {
}

func (s *Default) Capture(ctx context.Context, err error) {
	fmt.Println("=========================")
	j, _ := json.Marshal(err)
	fmt.Println(string(j))
	fmt.Println("=========================")
	sentry.CaptureException(err)
}
