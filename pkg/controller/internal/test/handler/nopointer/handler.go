package nopointer

import (
	"context"

	"github.com/giantswarm/operatorkit/v7/pkg/handler"
)

type Handler struct{}

func (h Handler) EnsureCreated(ctx context.Context, req handler.Request) (*handler.Response, error) {
	return nil, nil
}

func (h Handler) EnsureDeleted(ctx context.Context, req handler.Request) (*handler.Response, error) {
	return nil, nil
}
